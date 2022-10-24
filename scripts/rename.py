#!/usr/bin/env python3

import argparse
import re
import os

parser = argparse.ArgumentParser(description='Process target architectures')
parser.add_argument('--target-arch', metavar='t', type=str,
                    help='the target architecture that is being used for build')

args = parser.parse_args()
'''
Currently goreleaser supports only v1 for amd64 and hence renaming the target folder to be like linux_amd64
'''
regex = r"(?P<arch>^.*)_(?P<version>v[0-9]*$)"
matches = re.finditer(regex, args.target_arch)
try:
	for matchNum, match in enumerate(matches,start=1):
		if len(match.groups()) > 1:
			if  "linux_amd64" == match.group(1) and "v1" == match.group(2): 
				src = os.sep.join([os.getcwd(),"dist",f"server_{match.group(1)}_{match.group(2)}"])
				dest =  os.sep.join([os.getcwd(),"dist",f"server_{match.group(1)}"])
				print(f"Renaming {src} to {dest} ")
				os.rename(src,dest)
except Exception as e:
	print(e)
	exit(1)
