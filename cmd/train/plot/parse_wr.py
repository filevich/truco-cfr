import argparse
import os
import re
import json

import sys
sys.path.append('cmd/_com')
import parse_utils

parser = argparse.ArgumentParser(description='List all .out files in a directory')
parser.add_argument('-d', '--directory', type=str, required=True, help='Directory to search for .out files')
parser.add_argument('-o', '--output', type=str, default=None, help='Path to output JSON file')
args = parser.parse_args()

data = {}

for root, dirs, files in os.walk(args.directory):
    for file in files:
        if file.endswith('.out'):
            file_path = os.path.join(root, file)
            with open(file_path, 'r') as f:
                lines = f.readlines()
                ale = []
                simple = []
                for line in lines:
                    match_wr =re.search(r'ale=([\d\.]+) .([\d\.]+)\.\.([\d\.]+).+?di=(\d+).+?det=([\d\.]+) .([\d\.]+)\.\.([\d\.]+).+?di=(\d+)', line)
                    if match_wr:
                        ale.append({
                            "wr": float(match_wr.group(1)),
                            "l": float(match_wr.group(2)),
                            "u": float(match_wr.group(3)),
                            "di": int(match_wr.group(4)),
                        })
                        simple.append({
                            "wr": float(match_wr.group(5)),
                            "l": float(match_wr.group(6)),
                            "u": float(match_wr.group(7)),
                            "di": int(match_wr.group(8)),
                        })
                data[file] = {
                    "ale": ale,
                    "simple": simple,
                }

output_path = os.path.join(args.directory, 'result.json') \
    if args.output is None else args.output

with open(output_path, 'w') as f:
    json.dump(data, f, indent=2)

print(f"\nResult saved to {output_path}")
