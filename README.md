# Listty

List based note taking app in a terminal interface

## Usage

Create a folder in your current directory called `text`, and create a file called `example.txt`. Then `go run listty`,
or run from executable.

## Keybindings

### Modes

#### Global

    Ctrl-Q          Close application
    Ctrl-S          Save tree to file
    Shift-Up/Down   Swap selection with item above/below, without changing its depth
    Tab             Indent item
    Shift-Tab       Outdent item

#### Select

(This is the **default** and **startup** mode for Listty)

    Enter           Enter Edit mode
    Shift-Enter     Create new item
    Delete          Delete item // future
    Up/Down         Select the item above/below current selection
    Left/Right      Collapse/Expand // future
    D               Duplicate item // future
    .               Intdent (alias)
    ,               Outdent (alias)
    Ctrl-Right      Limit scope to selected item and its sub-tree // future
    Ctrl-Left       Increase scope to one level above selected item //future

#### Edit

(Type to add text to an item)

    Enter           Save item changes, leave Edit mode
    Shift-Enter     Save item changes, create and select a new item, in Edit mode
    Escape          Cancel item edit, return to previous state and enter selection mode
    Up/Down         Navigate to line begin/end
    Left/Right      Move cursor left/right
    Backspace       Remove character left of cursor
    
    
