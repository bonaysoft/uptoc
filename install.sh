#!/bin/sh
version="1.1"
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[1;34m'
YELLOW='\033[1;33m'
DARK='\033[1;30m'
NC='\033[0m'

echo "${BLUE}Uptoc binary installer ${version}${NC}"
unameOut="$(uname -s)"

case "${unameOut}" in
    Darwin*)    arch=macos
    			bin_dir="/usr/local/bin"
    			;;
    *)          arch=amd64
    			bin_dir="${HOME}/bin"
    			;;
esac
url=$(curl -s https://api.github.com/repos/saltbo/uptoc/releases/latest | grep "browser_download_url.*${arch}.tar.gz\"" | cut -d : -f 2,3 | tr -d '\"[:space:]')

echo "${DARK}"
echo "Configuration: [${arch}]"
echo "Location:      [${url}]"
echo "Directory:     [${bin_dir}]"
echo "${NC}"

test ! -d ${bin_dir} && mkdir ${bin_dir}
curl -J -L ${url} | tar xz -C ${bin_dir}

if [ $? -eq 0 ]
then
  echo "${GREEN}"
  echo "Installation completed successfully."
#  echo "$ wallet713 --version"
#  ${bin_dir}/wallet713 --version
#  if echo ":$PATH:"| grep -q  ".*:${bin_dir}:.*"; then
#    echo "\nRun 'wallet713' from the command line to open the wallet"
#  else
#  	echo "${YELLOW}\nPlease open a new terminal window to run 'wallet713' from the command-line"
#  fi
else
  echo "${RED}"
  echo "Failed installing uptoc"
fi

echo "${NC}"