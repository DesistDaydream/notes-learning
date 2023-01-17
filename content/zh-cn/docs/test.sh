#!/bin/bash

# specify the parent folder to search for subfolders
parent_folder="./IT学习笔记"

# search for subfolders in the parent folder
find ${parent_folder} -type d -print0 | while IFS= read -r -d '' dir; do
    # check if _index.md file already exists in the subfolder
    if [ ! -f "${dir}/_index.md" ]; then
        # create _index.md file and write text to it
        echo "---" > "${dir}/_index.md"
        echo "title: $(basename "${dir}")" >> "${dir}/_index.md"
        echo "---" >> "${dir}/_index.md"
    fi
done