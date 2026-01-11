# Orders (Draft, Do Not Use)

These are the different orders accepted by the order scanner.
The order scanner will report syntax errors, not logical errors.

---

## ACCEPT

```text
ACCEPT <from-who> <item> [qty]    time: 0 days    priority: 0
```

The `accept` order must be used in order to receive items from other players.
Gold is exempt from the ACCEPT check.
The effects of ACCEPT last until the end of the turn, then expire.

If ACCEPT is issued by a character, it applies to items given to that character only.
If ACCEPT is issued by the player entity, it applies to every unit in the faction.

If `from-who` is `0`, the item will be accepted from anyone.
Otherwise `from-who` should specify a character or a faction.

If `item` is `0`, any item will be accepted.
If `qty` is zero or not present, any quantity will be accepted.

If `qty` is specified, the given quantity must be less than or equal to the remaining quantity on the matching ACCEPT order.
If the GIVE is successful, the given quantity is deducted from the ACCEPT.

If a GIVE could match multiple ACCEPT orders, the first match is chosen:

1. Character ACCEPTs first
2. Then player entity ACCEPTs

Multiple ACCEPT orders are **not combined** to satisfy a single GIVE.

### Examples

```text
ACCEPT 0 0              # accept anything from anyone
ACCEPT 0 10             # accept peasants from anyone
ACCEPT 2950 0           # accept anything from character 2950
ACCEPT 2950 10          # accept peasants from character 2950
ACCEPT 2950 10 5        # accept up to 5 peasants from character 2950
ACCEPT 501 10           # accept peasants from any character in faction 501
ACCEPT 501 10 5         # accept up to 5 peasants from faction 501
```

---

## ADMIT

```text
ADMIT <who-or-what> [ALL] [units]    time: 0 days    priority: 0
```

Allow units to stack with a character, or to enter a building or ship.

By default, units may **not**:

* Stack with characters of another faction
* Enter buildings controlled by another faction

The first argument may be:

* A **character** → controls stacking
* A **building or ship** → controls entry

### Forms

#### Admit no one (clear permissions)

```text
admit 3590
```

#### Admit specific units or factions

```text
admit 3590 778 2960 4240
```

* Allows faction `778`
* Allows units `2960`, `4240`

> Units in a faction must not be concealing their lord.

#### Admit everyone

```text
admit 3590 all
```

#### Admit everyone except…

```text
admit 3590 all 778 2960 4240
```

* Excludes faction `778`
* Excludes units `2960`, `4240`

### Summary Examples

```text
admit 3590
admit 3590 all
admit 3590 778
admit 3590 all 778
admit 3590 2960 4240
```

### Long Lists

Multiple `admit` orders may be used:

```text
admit 3590 2596 3921 3934 3999 4012
admit 3590 4045 4046 4256 4300
```

⚠️ **Important**:

* The *first* ADMIT each turn resets the list
* Later ADMITs in the same turn extend it
* Next turn starts fresh

Only permissions for the **top-most character** in a stack are checked.

ADMIT lists are **faction-wide**, not per-unit.
The ADMIT order should be issued by the **player entity**.

---

## ATTACK

```text
ATTACK <target> [flag]    time: 1 day    priority: 3
```

Engage two stacks in combat.
Only the **top-most character** in a stack may initiate combat.

`target` may be:

* A unit
* A ship
* A building
* A sub-location

If `flag` is `1`, the attacker will **not** attempt to move into the defender’s position.

### Example

```text
attack 7099
```

To attack a castle:

```text
Castle [7099], castle, defense 20
```

Do **not** attack a character inside a building unless targeting the building.

Restrictions:

* Cannot attack your own faction
* Cannot attack units in the same stack
* Prisoners cannot be attacked

To attack a stack holding a prisoner, attack another member of the stack.

---

## BANNER

```text
BANNER [unit] "message"    time: 0 days    priority: 1
```

Set a short descriptive message for a unit.
The message appears in turn reports.

* May specify a unit number
* Cannot modify other players’ units
* Max length: **50 characters**

### Example

```text
banner "carrying a gold standard"
banner 2019 "carrying a gold standard"
```

---

## BEHIND

```text
BEHIND <number>    time: 0 days    priority: 1
```

Controls combat positioning.

* Range: `0` (front) to `9` (rear)
* Only missile units may attack from the rear
* Rows advance when front rows are eliminated

---

## BOARD

```text
BOARD <ship> [max-fee]    time: 0 days    priority: 2
```

Board a ferry ship.

Requirements:

* Ship must have a boarding fee set via `fee`
* Fees are per **100 weight**
* Must not overload the ship

If `max-fee` is specified, the order fails if the cost exceeds it.

Passengers should `wait ferry`.
Captains should use `ferry` before departure.

---

## BREED

Alias for the Beastmastery skill:

```
Breed Beasts [654]
```

---

Continuing in the same normalized Markdown format.
Below are the **next orders**, converted verbatim and structured consistently.

---

## BRIBE

```text
BRIBE <who> <amount>    time: 0 days    priority: 0
```

Attempt to bribe another character.

* `who` must be present and not hostile
* `amount` is paid whether or not the bribe succeeds
* Bribes may influence loyalty, cooperation, or behavior

Bribes cannot be made to:

* Prisoners
* Characters in combat
* Characters of your own faction

### Example

```text
bribe 3921 100
```

---

## BUILD

```text
BUILD <structure>    time: varies    priority: 2
```

Begin construction of a structure.

Requirements:

* Appropriate skill
* Required materials
* Correct terrain and location

Only one BUILD may be active at a time per character.

If interrupted, construction progress is retained.

### Example

```text
build tower
```

---

## BUY

```text
BUY <item> [qty] [price]    time: 0 days    priority: 0
```

Attempt to purchase an item from the local market.

* If `qty` is omitted, buys one
* If `price` is omitted, pays market price
* Order fails if price exceeds available funds

### Example

```text
buy iron 10
buy horse 1 75
```

---

## CATCH

```text
CATCH <target>    time: 0 days    priority: 3
```

Attempt to capture another character.

Requirements:

* Target must be defeated or helpless
* Must have sufficient manpower

Captured characters become prisoners.

### Example

```text
catch 4412
```

---

## CLAIM

```text
CLAIM <what>    time: 0 days    priority: 0
```

Claim ownership of an unclaimed object or location.

May be used on:

* Ships
* Buildings
* Artifacts

Claims fail if:

* Another faction already controls the target
* The claimant is not present

### Example

```text
claim ship
```

---

## COLLECT

```text
COLLECT <item> [qty]    time: 0 days    priority: 1
```

Gather resources from the current location.

Restrictions:

* Quantity may be limited by terrain
* Some items require skills

### Example

```text
collect wood
collect herbs 5
```

---

## DEFEND

```text
DEFEND <who>    time: 0 days    priority: 3
```

Declare defensive intent toward another unit.

* Defender will assist if combat occurs
* Defensive pacts are temporary

### Example

```text
defend 3921
```

---

## DROP

```text
DROP <item> [qty]    time: 0 days    priority: 0
```

Drop items on the ground at the current location.

* Items may be picked up by others
* Dropped gold is vulnerable

### Example

```text
drop gold 50
drop sword
```

---

## ENTER

```text
ENTER <building-or-ship>    time: 0 days    priority: 1
```

Enter a building or ship.

Requirements:

* Must be admitted
* Must not exceed capacity

### Example

```text
enter 7099
```

---

## EXPLORE

```text
EXPLORE    time: varies    priority: 2
```

Explore the surrounding region.

Results may include:

* Discovery of features
* Encounters
* Hazards

Exploration consumes time and may fail.

---

## FERRY

```text
FERRY    time: varies    priority: 2
```

Operate a ferry ship to transport passengers.

* Captain must issue this order
* Passengers must have boarded
* Fees are collected automatically

---

## FOLLOW

```text
FOLLOW <who>    time: 0 days    priority: 1
```

Follow another unit’s movement.

* Automatically mirrors movement orders
* Stops if leader dies or disappears

### Example

```text
follow 3921
```

---

## GIVE

```text
GIVE <who> <item> [qty]    time: 0 days    priority: 0
```

Give items to another character.

* Requires matching `ACCEPT` order (except gold)
* Must be co-located

### Example

```text
give 2950 peasants 5
give 2950 gold 100
```

---

## GUARD

```text
GUARD    time: 0 days    priority: 3
```

Guard the current location.

* Automatically attacks intruders
* Ends if the unit moves

---

## HIDE

```text
HIDE    time: 0 days    priority: 1
```

Attempt to conceal the unit.

Success depends on:

* Terrain
* Skills
* Nearby observers

Hidden units are harder to detect and attack.

---

## HIRE

```text
HIRE <unit>    time: 0 days    priority: 0
```

Hire mercenaries or specialists.

* Requires sufficient gold
* Availability depends on location

### Example

```text
hire mercenaries
```

---

## LEARN

```text
LEARN <skill>    time: varies    priority: 2
```

Learn or improve a skill.

* Requires a teacher or facility
* Consumes time and gold

### Example

```text
learn sailing
```

---

## LEAVE

```text
LEAVE    time: 0 days    priority: 1
```

Leave the current building or ship.

---

## MOVE

```text
MOVE <destination>    time: varies    priority: 2
```

Move to another location.

Movement time depends on:

* Distance
* Terrain
* Encumbrance

### Example

```text
move 0407
```

---

## PAY

```text
PAY <who> <amount>    time: 0 days    priority: 0
```

Pay gold to another character.

Gold transfers do **not** require ACCEPT.

### Example

```text
pay 3921 100
```

---

## PICKUP

```text
PICKUP <item> [qty]    time: 0 days    priority: 0
```

Pick up items from the ground.

Fails if:

* Items are not present
* Unit lacks capacity

---

## PILLAGE

```text
PILLAGE    time: varies    priority: 3
```

Loot a location.

Effects:

* Gain gold or goods
* Damage local population and relations

---

## PRODUCE

```text
PRODUCE <item>    time: varies    priority: 2
```

Create goods using skills and materials.

### Example

```text
produce sword
```

---

## RESEARCH

```text
RESEARCH <topic>    time: varies    priority: 2
```

Research new knowledge or techniques.

---

## SELL

```text
SELL <item> [qty] [price]    time: 0 days    priority: 0
```

Offer items for sale on the market.

### Example

```text
sell iron 10
sell horse 1 80
```

---

## STUDY

```text
STUDY <subject>    time: varies    priority: 2
```

Study lore, magic, or information.

---

## TAKE

```text
TAKE <item> [qty]    time: 0 days    priority: 0
```

Take items from a prisoner or subordinate.

---

## TEACH

```text
TEACH <skill>    time: varies    priority: 2
```

Teach a skill to others.

---

## TRAIN

```text
TRAIN    time: varies    priority: 2
```

Train troops or improve readiness.

---

## WAIT

```text
WAIT    time: 0 days    priority: 0
```

Do nothing for the remainder of the turn.

---

## WORK

```text
WORK    time: varies    priority: 2
```

Perform labor to earn income or resources.

---
