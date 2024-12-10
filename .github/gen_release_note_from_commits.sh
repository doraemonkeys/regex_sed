#!/bin/bash

while getopts "v:" opt; do
    case $opt in
    v)
        version_range=$OPTARG
        ;;
    \?)
        echo "Invalid option: -$OPTARG" >&2
        exit 1
        ;;
    esac
done

echo "version_range: $version_range"

if [ -z "$version_range" ]; then
    echo "Please provide the version range using -v option. Example: ./gen_release_note_from_commits.sh -v v1.14.1...v1.14.2"
    exit 1
fi

if [[ "$version_range" == ...* ]]; then
    version_range=${version_range#...}
fi

{
    echo "## What's Changed"
    git log --pretty=format:"* %h %s by %an@" --grep="^feat" -i "$version_range" | sort -f | uniq
    echo ""
    echo "## BUG & Fix"
    git log --pretty=format:"* %h %s by %an@" --grep="^fix" -i "$version_range" | sort -f | uniq
    echo ""
    echo "## Maintenance"
    git log --pretty=format:"* %h %s by %an@" --grep="^chore\|^docs\|^refactor" -i "$version_range" | sort -f | uniq
    echo ""
    # substitute your username and repo
    echo "**Full Changelog**: https://github.com/doraemonkeys/sedr/compare/$version_range"
} >release.md
