# Created at 01 11 2020
# Created by miklos.lorinczi@appventurez.nl

# This script needs to be sourced by Bitbucket Pipeline' steps
# script:
#   - source scripts/addenv.sh

# Then variables can be saved with: addenv <key> <value>

# The script needs to be saved as an artifact for further steps:
# artifacts:
#   - scripts/addenv.sh

# Use this function in your step to add variables to this file
function addenv() {
    echo "export ${1}=${2}" >> scripts/addenv.sh
}
