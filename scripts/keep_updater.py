import os
import gkeepapi
import json
import sys


def get_node(title, keep):
    gnotes = keep.find(func=lambda x: x.title == title)
    if gnotes:
        for g in gnotes:
            if g.title == title:
                return g
    print(f"Note with title '{title}' not found.")
    return None


def get_list_item(note, item_text):
    items = note.items
    for i in items:
        if i.text == item_text:
            return i
    return None


def clear_items(item):
    subitems = item.subitems
    for sub in subitems:
        print(f"Clearing subitem: {sub.text}")
        sub.delete()


def add_items(note, parent_note, ingredients):
    if not ingredients:
        print("No ingredients to add.")
        return
    for ingredient in ingredients:
        item_text = f"{ingredient['name']}   ({ingredient['quantity']})"
        print(f"Adding item: {item_text}")
        new_note = note.add(item_text, False)
        parent_note.indent(new_note)
    print("All items added.")


def update_note(note, ingredients):
    print(f"Updating note: {note.title}")
    parent_item = get_list_item(note, "Zenful Shopping Items")
    if not parent_item:
        print(f"Creating parent item: Zenful Shopping Items")
        parent_item = note.add("Zenful Shopping Items", False)
    print(f"Clearing existing items")
    clear_items(parent_item)
    if ingredients is None:
        parent_item.checked = True
        print(f"No ingredients provided for {note.title}. Skipping update.")
    else:
        parent_item.checked = False
    print(f"Adding new items")
    add_items(note, parent_item, ingredients)
    print(f"{note.title} updated succesfully")


def sync(data):
    username = os.getenv("GOOGLE_USERNAME")
    master_token = os.getenv("GOOGLE_TOKEN")
    if not master_token:
        print(
            "No master token found. Please set the GOOGLE_TOKEN environment variable."
        )
        return
    if not username:
        print("No username found. Please set the GOOGLE_USERNAME environment variable.")
        return

    keep = gkeepapi.Keep()
    keep.authenticate(username, master_token)

    for l in data:
        note = get_node(l["title"], keep)
        if not note:
            print(f"Creating new note: {l['title']}")
            note = keep.createList(l["title"], [])
        else:
            print(f"Found existing note: {note.title}")
        note.color = gkeepapi.node.ColorValue(l["color"])
        print(f"Setting note color to {note.color}")
        update_note(note, l["ingredients"])

    keep.sync()


def main():
    """
    Reads JSON data from stdin, which should contain a list of notes with their titles,
    ingredients, and colors.
    The expected JSON format is:
    [
        {
            "title": "Note Title",
            "ingredients": [
                {"name": "Ingredient Name", "quantity": "Quantity"},
                ...
            ],
            "color": "Color Name"
        },
        ...
    ]
    """

    data = json.load(sys.stdin)
    if not data:
        print("No data provided.")
        return
    sync(data)


if __name__ == "__main__":
    main()
