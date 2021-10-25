#!/bin/bash

# Leverage the default env variables as described in:
# https://docs.github.com/en/actions/reference/environment-variables#default-environment-variables
if [[ $GITHUB_ACTIONS != "true" ]]
then
  /usr/bin/yor $@
  exit $?
fi

flags=""

# Actions pass inputs as $INPUT_<input name> environment variables
[[ -n "$INPUT_TAG_GROUPS" ]] && flags="$flags--tag-groups $INPUT_TAG_GROUPS "
[[ -n "$INPUT_TAG" ]] && flags="$flags--tag $INPUT_TAG "
[[ -n "$INPUT_SKIP_TAGS" ]] && flags="$flags--skip-tags $INPUT_SKIP_TAGS "
[[ -n "$INPUT_SKIP_DIRS" ]] && flags="$flags--skip-dirs $INPUT_SKIP_DIRS "
[[ -n "$INPUT_SKIP_RESOURCE_TYPES" ]] && flags="$flags--skip-resource-types $INPUT_SKIP_RESOURCE_TYPES "
[[ -n "$INPUT_CUSTOM_TAGS" ]] && flags="$flags--custom-tagging $INPUT_CUSTOM_TAGS "
[[ -n "$INPUT_OUTPUT_FORMAT" ]] && flags="$flags--output $INPUT_OUTPUT_FORMAT "
[[ -n "$INPUT_CONFIG_FILE" ]] && flags="$flags--config-file $INPUT_CONFIG_FILE "
[[ -n "$INPUT_LOG_LEVEL" ]] && export LOG_LEVEL=$INPUT_LOG_LEVEL

[[ -d ".yor_plugins" ]] && echo "Directory .yor_plugins exists, and will be overwritten by yor. Please rename this directory."

echo "running command:"
echo yor tag -d $INPUT_DIRECTORY $flags

/usr/bin/yor tag -d $INPUT_DIRECTORY $flags
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
