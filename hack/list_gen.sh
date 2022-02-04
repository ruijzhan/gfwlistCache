#!/bin/bash

wget https://raw.githubusercontent.com/ruijzhan/chnroute/master/gfwlist.txt -O tmp1
sed -i 's/^/gfwlist /g' tmp1

wget https://raw.githubusercontent.com/privacy-protection-tools/anti-AD/master/anti-ad-domains.txt -O tmp2
sed -i 's/^/adblock /g' tmp2


cat tmp* > list.txt
rm tmp*
