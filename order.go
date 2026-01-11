// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

import (
	"database/sql"
	"fmt"
)

// order.c -- manage list of unit orders for each faction
//
// This file ports order.c from the C codebase. In the Go version,
// orders are loaded from the database rather than parsed from files.
// Text parsing is skipped since the browser submits structured orders
// directly to the DB.

// OrderEntry represents a single order for a unit.
// In the C code, this was a char* in order_list.l plist.
type OrderEntry struct {
	ID           int    // Database ID (from orders table)
	TurnNumber   int    // Turn this order is for
	PlayerID     int    // Player who submitted the order
	UnitID       int    // Unit (character) the order is for
	RawText      string // The raw order text
	SourceChannel string // Where the order came from (email, web, etc.)
}

// OrderQueue holds the pending orders for a specific unit.
// This replaces the C order_list struct.
type OrderQueue struct {
	Unit   int           // unit orders are for
	Orders []*OrderEntry // list of orders for unit (replaces **char l)
}

// p_order_head returns the order queue for a unit, creating one if needed.
// Port of C p_order_head().
func (e *Engine) p_order_head(pl, who int) *OrderQueue {
	p := e.rp_player(pl)
	if p == nil {
		return nil
	}

	// Check if we already have an order queue for this unit
	if e.globals.orderQueues == nil {
		e.globals.orderQueues = make(map[int]map[int]*OrderQueue)
	}
	if e.globals.orderQueues[pl] == nil {
		e.globals.orderQueues[pl] = make(map[int]*OrderQueue)
	}

	if q, ok := e.globals.orderQueues[pl][who]; ok {
		return q
	}

	// Create new order queue
	q := &OrderQueue{Unit: who}
	e.globals.orderQueues[pl][who] = q
	return q
}

// rp_order_head returns the order queue for a unit, or nil if none exists.
// Port of C rp_order_head().
func (e *Engine) rp_order_head(pl, who int) *OrderQueue {
	if e.globals.orderQueues == nil {
		return nil
	}
	if e.globals.orderQueues[pl] == nil {
		return nil
	}
	return e.globals.orderQueues[pl][who]
}

// top_order returns the first order in the queue for a unit, or empty string if none.
// Port of C top_order().
func (e *Engine) top_order(pl, who int) string {
	q := e.rp_order_head(pl, who)
	if q == nil || len(q.Orders) == 0 {
		return ""
	}
	return q.Orders[0].RawText
}

// is_stop_order checks if an order string is a "stop" command.
// Port of C is_stop_order().
func is_stop_order(s string) bool {
	if s == "" {
		return false
	}

	// Skip leading whitespace
	i := 0
	for i < len(s) && iswhite(s[i]) {
		i++
	}

	// Extract first word
	j := i
	for j < len(s) && !iswhite(s[j]) {
		j++
	}
	word := s[i:j]

	if i_strcmp(word, "stop") == 0 {
		return true
	}

	if fuzzy_strcmp(s, "stop") {
		return true
	}

	return false
}

// stop_order returns true if a stop order is queued for the given unit.
// STOP orders must be the first command in the order queue.
// Port of C stop_order().
func (e *Engine) stop_order(pl, who int) bool {
	s := e.top_order(pl, who)
	if s == "" {
		return false
	}
	return is_stop_order(s)
}

// pop_order removes the first order from a unit's queue.
// Port of C pop_order().
func (e *Engine) pop_order(pl, who int) {
	q := e.rp_order_head(pl, who)
	if q == nil || len(q.Orders) == 0 {
		return
	}
	q.Orders = q.Orders[1:]
}

// flush_unit_orders removes all orders for a unit.
// Port of C flush_unit_orders().
func (e *Engine) flush_unit_orders(pl, who int) {
	if e.globals.bx[who] == nil {
		return
	}

	for e.top_order(pl, who) != "" {
		e.pop_order(pl, who)
	}

	if e.player(who) == pl {
		c := e.rp_command(who)
		if c != nil && c.state == STATE_LOAD {
			e.command_done(c)
		}
	}
}

// queue_order appends an order to a unit's queue.
// Port of C queue_order().
func (e *Engine) queue_order(pl, who int, s string) {
	q := e.p_order_head(pl, who)
	if q == nil {
		return
	}

	// Limit queue size to 250 orders (same as C)
	if len(q.Orders) >= 250 {
		return
	}

	q.Orders = append(q.Orders, &OrderEntry{
		PlayerID: pl,
		UnitID:   who,
		RawText:  s,
	})
}

// prepend_order adds an order to the front of a unit's queue.
// Port of C prepend_order().
func (e *Engine) prepend_order(pl, who int, s string) {
	q := e.p_order_head(pl, who)
	if q == nil {
		return
	}

	entry := &OrderEntry{
		PlayerID: pl,
		UnitID:   who,
		RawText:  s,
	}
	q.Orders = append([]*OrderEntry{entry}, q.Orders...)
}

// queue_stop prepends a stop order if one isn't already queued.
// Port of C queue_stop().
func (e *Engine) queue_stop(pl, who int) {
	if e.stop_order(pl, who) {
		return
	}
	e.prepend_order(pl, who, "stop")
}

// queue is a convenience wrapper for queue_order that uses the unit's player.
// Port of C queue() - but without the varargs formatting since Go has fmt.Sprintf.
func (e *Engine) queue(who int, format string, args ...any) {
	pl := e.player(who)
	s := fmt.Sprintf(format, args...)
	e.queue_order(pl, who, s)
}

// LoadOrders loads orders from the database for the current turn.
// This replaces the C load_orders() function which read from files.
func (e *Engine) LoadOrders(turnNumber int) error {
	if e.globals.orderQueues == nil {
		e.globals.orderQueues = make(map[int]map[int]*OrderQueue)
	}

	rows, err := e.db.Query(`
		SELECT id, turn_number, player_id, source_char_id, raw_text, source_channel
		FROM orders
		WHERE turn_number = ?
		ORDER BY id
	`, turnNumber)
	if err != nil {
		return fmt.Errorf("query orders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, turn, playerID int
		var unitID sql.NullInt64
		var rawText string
		var sourceChannel sql.NullString

		if err := rows.Scan(&id, &turn, &playerID, &unitID, &rawText, &sourceChannel); err != nil {
			return fmt.Errorf("scan order %d: %w", id, err)
		}

		// Skip orders without a unit
		if !unitID.Valid {
			continue
		}

		who := int(unitID.Int64)

		// Get or create order queue for this player/unit
		if e.globals.orderQueues[playerID] == nil {
			e.globals.orderQueues[playerID] = make(map[int]*OrderQueue)
		}
		if e.globals.orderQueues[playerID][who] == nil {
			e.globals.orderQueues[playerID][who] = &OrderQueue{Unit: who}
		}

		entry := &OrderEntry{
			ID:           id,
			TurnNumber:   turn,
			PlayerID:     playerID,
			UnitID:       who,
			RawText:      rawText,
		}
		if sourceChannel.Valid {
			entry.SourceChannel = sourceChannel.String
		}

		e.globals.orderQueues[playerID][who].Orders = append(
			e.globals.orderQueues[playerID][who].Orders, entry)
	}

	return rows.Err()
}

// SaveOrders saves all pending orders to the database.
// This replaces the C save_orders() function which wrote to files.
func (e *Engine) SaveOrders(turnNumber int) error {
	// Delete existing orders for this turn first
	_, err := e.db.Exec(`DELETE FROM orders WHERE turn_number = ?`, turnNumber)
	if err != nil {
		return fmt.Errorf("delete old orders: %w", err)
	}

	// Insert all orders from memory
	for playerID, unitQueues := range e.globals.orderQueues {
		for _, queue := range unitQueues {
			if queue == nil || len(queue.Orders) == 0 {
				continue
			}

			// Skip invalid units
			if e.globals.bx[queue.Unit] == nil || e.Kind(queue.Unit) == T_deadchar {
				continue
			}

			for _, order := range queue.Orders {
				var sourceChannel sql.NullString
				if order.SourceChannel != "" {
					sourceChannel = sql.NullString{String: order.SourceChannel, Valid: true}
				}

				_, err := e.db.Exec(`
					INSERT INTO orders (turn_number, player_id, source_char_id, raw_text, source_channel)
					VALUES (?, ?, ?, ?, ?)
				`, turnNumber, playerID, queue.Unit, order.RawText, sourceChannel)
				if err != nil {
					return fmt.Errorf("insert order for unit %d: %w", queue.Unit, err)
				}
			}
		}
	}

	return nil
}

// ClearOrders removes all in-memory order queues.
func (e *Engine) ClearOrders() {
	e.globals.orderQueues = nil
}

// GetOrderQueue returns the order queue for a player/unit pair.
// This is a read-only accessor for testing and inspection.
func (e *Engine) GetOrderQueue(playerID, unitID int) *OrderQueue {
	return e.rp_order_head(playerID, unitID)
}

// CountOrders returns the number of orders queued for a unit.
func (e *Engine) CountOrders(playerID, unitID int) int {
	q := e.rp_order_head(playerID, unitID)
	if q == nil {
		return 0
	}
	return len(q.Orders)
}

// GetAllOrders returns all orders for a unit as a slice of strings.
func (e *Engine) GetAllOrders(playerID, unitID int) []string {
	q := e.rp_order_head(playerID, unitID)
	if q == nil {
		return nil
	}
	result := make([]string, len(q.Orders))
	for i, o := range q.Orders {
		result[i] = o.RawText
	}
	return result
}

// helper stubs for functions that will be implemented in other sprints

// rp_command returns the command structure for an entity, or nil.
// This is a stub that will be properly implemented when command execution is ported.
func (e *Engine) rp_command(who int) *command {
	if e.globals.bx[who] == nil {
		return nil
	}
	return e.globals.bx[who].cmd
}

// command_done marks a command as completed.
// This is a stub that will be properly implemented when command execution is ported.
func (e *Engine) command_done(c *command) {
	if c != nil {
		c.state = STATE_DONE
	}
}

// player returns the player ID that owns a unit.
// This is a stub - the full implementation exists in accessor.go or similar.
func (e *Engine) player(who int) int {
	if e.globals.bx[who] == nil {
		return 0
	}
	ch := e.globals.bx[who].x_char
	if ch == nil {
		return 0
	}
	return ch.unit_lord
}

// rp_player returns the entity_player for a player, or nil.
func (e *Engine) rp_player(pl int) *entity_player {
	if e.globals.bx[pl] == nil {
		return nil
	}
	return e.globals.bx[pl].x_player
}
