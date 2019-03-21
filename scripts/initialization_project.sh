#!/bin/bash
PROJECT_NAME=$1

if [ "$PROJECT_NAME" == "" ]; then
    echo "not found project_name"
    echo " * ex: ./initialization_project.sh project_name"
    exit
fi

PROJECT_NAME_LOWER=$(echo $PROJECT_NAME | tr '[:upper:]' '[:lower:]')

git grep -l 'magneto' | xargs sed -i '' -e "s/magneto/$PROJECT_NAME/g" | sh
find ./ -type d -name '*magneto*' | awk '(so=$0){gsub(/magneto/,"'$PROJECT_NAME'",$0);print "mv " so " " $0}' | sh
find ./  -name '*magneto*' | awk '(so=$0){gsub(/magneto/,"'$PROJECT_NAME'",$0);print "mv " so " " $0}' | sh

echo "... initialization done"
rm -rf .git ../scripts/initialization_project.sh
