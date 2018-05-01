#!/bin/bash

dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cli="${dir}/../cli/bin/kv"
store_name="test-int-store"
key_name="test-key-1"

function cleanup() {
  echo "Cleaning up"
  ${cli} delete-store --store-name ${store_name}
  exit $exit_code
}

# successful store creation
output=$(${cli} create-store --store-name ${store_name} 2>&1)
status=$?
if [ $status -ne 0 ]; then
  echo "create-store should be successful but wasn't (code: $? | output: ${output})"
  cleanup 1
fi

# unsuccessful duplicate store
output=$(${cli} create-store --store-name ${store_name} 2>&1)
status=$?
if [ $status -eq 0 ]; then
  echo "create-store should not be successful but was (code: $? | output: ${output})"
  cleanup 1
fi

# successful key set
output=$(${cli} set-key --store-name ${store_name} --key-name ${key_name} --key-value hi 2>&1)
status=$?
if [ $status -ne 0 ]; then
  echo "set-key should be successful but wasn't (code: $? | output: ${output})"
  cleanup 1
fi

# successful get key
output=$(${cli} get-key --store-name ${store_name} --key-name ${key_name} 2>&1)
status=$?
if [ $status -ne 0 ] || [ "${output}" != "hi" ]; then
  echo "get-key should be successful but wasn't (code: $? | output: ${output})"
  cleanup 1
fi

# successful delete key
output=$(${cli} delete-key --store-name ${store_name} --key-name ${key_name} 2>&1)
status=$?
if [ $status -ne 0 ]; then
  echo "delete-key should be successful but wasn't (code: $? | output: ${output})"
  cleanup 1
fi

# successful delete non-existent key
output=$(${cli} delete-key --store-name ${store_name} --key-name bad 2>&1)
status=$?
if [ $status -ne 0 ]; then
  echo "delete-key should be successful but wasn't (code: $? | output: ${output})"
  cleanup 1
fi

# successful delete store
output=$(${cli} delete-store --store-name ${store_name} 2>&1)
status=$?
if [ $status -ne 0 ]; then
  echo "delete-store should be successful but wasn't (code: $? | output: ${output})"
  cleanup 1
fi

# successful delete non-existent store
output=$(${cli} delete-store --store-name bad 2>&1)
status=$?
if [ $status -ne 0 ]; then
  echo "delete-store should be successful but wasn't (code: $? | output: ${output})"
  cleanup 1
fi

echo "All tests successful"
cleanup 0
