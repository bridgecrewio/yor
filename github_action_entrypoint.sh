#!/bin/bash

# Actions pass inputs as $INPUT_<input name> environmet variables

[[ -n "$INPUT_TAG_GROUPS" ]] && TAG_GROUPS="--tag-groups $INPUT_TAG_GROUPS"
[[ -n "$INPUT_TAG" ]] && TAG_FLAG="--tag $INPUT_TAG"
[[ -n "$INPUT_SKIP_TAGS" ]] && SKIP_TAG_FLAG="--skip-tags $INPUT_SKIP_TAGS"
[[ -n "$INPUT_SKIP_DIRS" ]] && SKIP_DIR_FLAG="--framework $INPUT_SKIP_TAGS"
[[ -n "$INPUT_CUSTOM_TAGS" ]] && EXT_TAGS_FLAG="--custom-tagging $INPUT_CUSTOM_TAGS"
[[ -n "$INPUT_OUTPUT_FORMAT" ]] && OUTPUT_FLAG="--output $INPUT_OUTPUT_FORMAT"
[[ -n "$INPUT_LOG_LEVEL" ]] && export LOG_LEVEL=$INPUT_LOG_LEVEL

[[ -d ".yor_plugins" ]] && echo "Directory .yor_plugins exists, and will be overwritten by yor. Please rename this directory."

echo "running yor on directory: $INPUT_DIRECTORY"
/usr/bin/yor tag -d "$INPUT_DIRECTORY" "$TAG_FLAG" "$TAG_GROUPS" "$SKIP_TAG_FLAG" "$SKIP_DIR_FLAG" "$EXT_TAGS_FLAG" "$OUTPUT_FLAG"
rm -rf .yor_plugins
YOR_EXIT_CODE=$?

_git_is_dirty() {
    [ -n "$(git status -s --untracked-files=no)" ]
}

if [[ $YOR_EXIT_CODE -eq 0 && $INPUT_COMMIT_CHANGES == "YES" ]]
then
  if _git_is_dirty
  then
    echo "Yor made changes, committing"
    git add .
    git -c user.name=actions@github.com -c user.email="GitHub Actions" \
        commit -m "Update tags (by Yor)" \
        --author="github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>" ;
    echo "Changes committed, pushing..."
    git push origin
  fi
else
  echo "::debug::exiting, yor failed or commit is skipped"
  echo "::debug::yor exit code: $YOR_EXIT_CODE"
  echo "::debug::commit_changes: $INPUT_COMMIT_CHANGES"
  exit $YOR_EXIT_CODE
fi