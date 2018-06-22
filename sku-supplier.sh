#!/bin/bash
# for supervisor use   sh sku-supplier.sh
#gopath=`echo $GOPATH`
#cd gopath

dirpath=`dirname $0`
cd $dirpath
echo `pwd`
echo $(date +%Y-%m-%d_%H:%M:%S)
nohup ./monitor &



