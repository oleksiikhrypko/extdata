#!/usr/bin/env bash
if [ -z "$1" ]
  then
    echo "Existing git tag required as a parameter"
    exit 1
fi

git checkout $1
if [ $? -eq 0 ];
then
    echo "Checkout was successful, version exists"
else
    echo "Version $1 does not exist in git. Try to pull latest code"
    exit 1
fi
echo "Cleaning up extra git tags..."
git tag | grep '^v.*-uat' | xargs -n 1 -r git tag -d
git fetch origin --prune-tags
git tag | grep '^v.*-uat' | xargs -n 1 -r git push --delete origin
git tag | grep '^v.*-uat' | xargs -n 1 -r git tag -d
echo "Releasing to uat..."
git tag -a "$1-uat" -m "$1-uat"
git push origin "$1-uat"
