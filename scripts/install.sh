#!/bin/sh

# version="$(cat plugin.yaml | grep "version" | cut -d '"' -f 2)"
# cd $HELM_PLUGIN_DIR
# version="$(cat plugin.yaml | grep "version" | cut -d '"' -f 2)"
# echo "Installing helm-bulk ${version} ..."
#
# # Find correct archive name
# unameOut="$(uname -s)"
#
# case "${unameOut}" in
#     Linux*)     os=Linux;;
#     Darwin*)    os=Darwin;;
#     CYGWIN*)    os=Cygwin;;
#     MINGW*)     os=windows;;
#     *)          os="UNKNOWN:${unameOut}"
# esac
#
# arch=`uname -m`
# url="https://github.com/ovotech/helm-bulk/releases/download/${version}/helm-bulk_${version}_${os}_${arch}.tar.gz"
#
# if [ "$url" = "" ]
# then
#     echo "Unsupported OS / architecture: ${os}_${arch}"
#     exit 1
# fi
#
# filename=`echo ${url} | sed -e "s/^.*\///g"`
#
# # Download archive
# if [ -n $(command -v curl) ]
# then
#     curl -sSL -O $url
# elif [ -n $(command -v wget) ]
# then
#     wget -q $url
# else
#     echo "Need curl or wget"
#     exit -1
# fi
#
# # Install bin
# rm -rf bin && mkdir bin && tar xzvf $filename -C bin > /dev/null && rm -f $filename

echo "helm-bulk ${version} is correctly installed."
