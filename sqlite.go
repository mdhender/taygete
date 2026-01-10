// Copyright (c) 2026 Michael D Henderson. All rights reserved.

package taygete

func (e *Engine) savePrngState(name string) error {
	state, err := e.prng.MarshalBinary()
	if err != nil {
		e.logger.Error("savePrngState: marshal failed", "name", name, "err", err)
		return err
	}
	_, err = e.db.Exec(`INSERT INTO prng_state (name, state) VALUES (?, ?) ON CONFLICT(name) DO UPDATE SET state = excluded.state `, name, state)
	if err != nil {
		e.logger.Error("savePrngState: update failed", "name", name, "err", err)
		return err
	}
	return nil
}

func (e *Engine) restorePrngState(name string) error {
	var state []byte
	err := e.db.QueryRow(`SELECT state FROM prng_state WHERE name = ?`, name).Scan(&state)
	if err != nil {
		e.logger.Error("restorePrngState: select failed", "name", name, "err", err)
		return err
	}
	err = e.prng.UnmarshalBinary(state)
	if err != nil {
		e.logger.Error("restorePrngState: unmarshal failed", "name", name, "err", err)
		return err
	}
	return nil
}

func (e *Engine) readPassword(key string) (string, error) {
	var value string
	err := e.db.QueryRow(`SELECT value FROM passwords WHERE key = ?`, key).Scan(&value)
	if err != nil {
		e.logger.Error("readPassword: select failed", "key", key, "err", err)
		return "", err
	}
	return value, nil
}

func (e *Engine) savePassword(key, value string) error {
	_, err := e.db.Exec(`INSERT INTO passwords (key, value) VALUES (?, ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value`, key, value)
	if err != nil {
		e.logger.Error("savePassword: save failed", "key", key, "err", err)
		return err
	}
	return nil
}
