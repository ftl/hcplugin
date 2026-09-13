# Hello Contest OpenAction Plugin

An OpenAction ([OpenDeck](https://github.com/nekename/OpenDeck)) plugin to control the Hello Contest logger.

## Configuration

The plugin currently supports only one the default connection to `localhost:8383`.

## Actions

The following actions are currently available:

| Action | Description |
|--------|-------------|
| Do Action | trigger one of the actions described in https://github.com/ftl/hellocontest/wiki/Configuration#keybindings |
| Send | send a morse code template as described in https://github.com/ftl/hellocontest/wiki/Configuration#keyer-presets |
| Search & Pounce | navigate through the spots and trigger the start of a QSO |

### Long Press

The key actions *Do Action* and *Send* also have a long press setting. A press of
500 ms or more uses that setting, a shorter press uses the normal setting. If the
long press setting is empty, a long press does the same as a short press. Both
actions run when you release the key. The dial press of an encoder always runs
the normal setting immediately.

## License
This software is published under the [MIT License](https://www.tldrlegal.com/l/mit).

Copyright [Florian Thienel](http://thecodingflow.com/)
