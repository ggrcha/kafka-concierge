#!bin/sh

L_GREEN='\033[1;32m'
RED='\033[0;31m'
NC='\033[0m'

PROJECT_NAME=$(pwd | awk '{n=split($0,a,"/"); print a[n]}')
PROJECT_GIT=$(git remote -v | head -n 1 | cut -d$' ' -f1 | cut -d$'\t' -f2)

echo "${L_GREEN}Project Name: $PROJECT_NAME"
echo "Project Git Repo: $PROJECT_GIT"
echo "${NC}"

if [ "$(curl -I -L https://api-semver.gtw.dev.io.bb.com.br/project/$PROJECT_NAME/info 2>/dev/null | head -n 1 | cut -d$' ' -f2)" -eq "403" ]; 
then 
    echo "Creating project at semver"
    curl -X POST https://api-semver.gtw.dev.io.bb.com.br/project/createProject -H 'Accept: application/json' -H 'Content-Type: application/json' -d '{
        "projectName": "'$PROJECT_NAME'",
        "codeBaseRepo": "'$PROJECT_GIT'",
        "releaseDevHook": "",
        "releaseHomHook": "",
        "releaseCanHook": "",
        "releaseProdHook": ""
    }';
fi

string=$(git describe --tags $(git rev-list --tags --max-count=1))
if [ -z "$string" ];
then
    string="0.0.0"
fi
array=(${string//./ })
if [ "$1" = "patch" ];
then
    echo "${L_GREEN}Increase Patch Version${NC}"
    array[2]=$((${array[2]}+1))
elif [ "$1" = "minor" ];
then
    echo "${L_GREEN}Increase Minor Version${NC}"
    array[1]=$((${array[1]}+1))
    array[2]=$((0))
elif [ "$1" = "major" ];
then 
    echo "${L_GREEN}Increase Major Version${NC}"    
    array[0]=$((${array[0]}+1))
    array[1]=$((0))
    array[2]=$((0))
else
    echo "${RED}Please inform one of this options (patch|minor|major) at command argument: sh tag.sh ${NC}"
fi

TAG=$(printf ".%s" "${array[@]}")
TAG=${TAG:1}
echo "Actual TAG: $string"
echo "New Tag: $TAG"

git tag $TAG
# curl -L -k -H 'Content-Type: application/json' -X POST -d '{ "projectName": "'$PROJECT_NAME'", "version":"'\"$(git describe --tags $(git rev-list --tags --max-count=1))\"'" }' api-semver.gtw.dev.io.bb.com.br/project/version/development
curl -L -k -H 'Content-Type: application/json' -X POST -d '{ "projectName": "'$PROJECT_NAME'", "version":"'$(git describe --tags $(git rev-list --tags --max-count=1))'" }' api-semver.gtw.dev.io.bb.com.br/project/version/development
git push
git push --tags
