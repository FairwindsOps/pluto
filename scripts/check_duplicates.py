#!/usr/bin/env python3
"""
Script to check for duplicate entries in versions.yaml file.
A duplicate is defined as having the same 'version' and 'kind' values.
"""

import sys
import yaml
from collections import defaultdict

def check_duplicates(yaml_file):
    """Check for duplicates in the versions.yaml file."""
    with open(yaml_file, 'r') as f:
        data = yaml.safe_load(f)
    
    # Track combinations of version and kind
    combinations = defaultdict(list)
    duplicates = []
    
    for i, item in enumerate(data['deprecated-versions']):
        version = item.get('version', '')
        kind = item.get('kind', '')
        combo = (version, kind)
        
        # Store the index where this combination appears
        combinations[combo].append(i)
    
    # Find duplicates
    for combo, indices in combinations.items():
        if len(indices) > 1:
            version, kind = combo
            duplicates.append({
                'version': version,
                'kind': kind,
                'indices': indices
            })
    
    if duplicates:
        print("ERROR: Duplicate entries found in versions.yaml:")
        for dup in duplicates:
            print(f"  - version: {dup['version']}, kind: {dup['kind']} (found at indices: {dup['indices']})")
        return 1
    
    print("No duplicates found in versions.yaml")
    return 0

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python check_duplicates.py <path_to_versions.yaml>")
        sys.exit(1)
    
    yaml_file = sys.argv[1]
    sys.exit(check_duplicates(yaml_file))