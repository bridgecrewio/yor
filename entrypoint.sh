#!/bin/bash

# Leverage the default env variables as described in:
# https://docs.github.com/en/actions/reference/environment-variables#default-environment-variables
if [[ $GITHUB_ACTIONS != "true" ]]
then
  /usr/bin/yor $@
  exit $?
fi

# Actions pass inputs as $INPUT_<input name> environment variables
[[ -n "$INPUT_TAG_GROUPS" ]] && TAG_GROUPS="--tag-groups $INPUT_TAG_GROUPS"
[[ -n "$INPUT_TAG" ]] && TAG_FLAG="--tag $INPUT_TAG"
[[ -n "$INPUT_SKIP_TAGS" ]] && SKIP_TAG_FLAG="--skip-tags $INPUT_SKIP_TAGS"
[[ -n "$INPUT_SKIP_DIRS" ]] && SKIP_DIR_FLAG="--skip-dirs $INPUT_SKIP_DIRS"
[[ -n "$INPUT_CUSTOM_TAGS" ]] && EXT_TAGS_FLAG="--custom-tagging $INPUT_CUSTOM_TAGS"
[[ -n "$INPUT_OUTPUT_FORMAT" ]] && OUTPUT_FLAG="--output $INPUT_OUTPUT_FORMAT"
[[ -n "$INPUT_CONFIG_FILE" ]] && CONFIG_FILE_FLAG="--config-file $INPUT_CONFIG_FILE"
[[ -n "$INPUT_LOG_LEVEL" ]] && export LOG_LEVEL=$INPUT_LOG_LEVEL

[[ -d ".yor_plugins" ]] && echo "Directory .yor_plugins exists, and will be overwritten by yor. Please rename this directory."

echo "running yor on directory: $INPUT_DIRECTORY"
echo "params"
echo INPUT_TAG_GROUPS "$INPUT_TAG_GROUPS"
echo INPUT_TAG "$INPUT_TAG"
echo INPUT_SKIP_TAGS "$INPUT_SKIP_TAGS"
echo INPUT_SKIP_DIRS "$INPUT_SKIP_DIRS"
echo INPUT_CUSTOM_TAGS "$INPUT_CUSTOM_TAGS"
echo INPUT_OUTPUT_FORMAT "$INPUT_OUTPUT_FORMAT"
echo INPUT_CONFIG_FILE "$INPUT_CONFIG_FILE"
echo INPUT_LOG_LEVEL "$INPUT_LOG_LEVEL"

echo "running command:"
echo yor tag -d "$INPUT_DIRECTORY" "$TAG_FLAG" "$TAG_GROUPS" "$SKIP_TAG_FLAG" "$SKIP_DIR_FLAG" "$EXT_TAGS_FLAG" "$OUTPUT_FLAG" "$CONFIG_FILE_FLAG"

/usr/bin/yor tag -d "$INPUT_DIRECTORY" "$TAG_FLAG" "$TAG_GROUPS" "$SKIP_TAG_FLAG" "$SKIP_DIR_FLAG" "$EXT_TAGS_FLAG" "$OUTPUT_FLAG" "$CONFIG_FILE_FLAG"
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
fi
exit $YOR_EXIT_CODE