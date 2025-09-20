#!/bin/sh
set -e

# Replace env placeholders inside .next static files
for file in $(find .next -type f -name '*.js'); do
  sed -i "s#NEXT_PUBLIC_API_URL#${NEXT_PUBLIC_API_URL}#g" $file
  sed -i "s#NEXT_PUBLIC_WS_URL#${NEXT_PUBLIC_WS_URL}#g" $file
  sed -i "s#NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY#${NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY}#g" $file
done

exec npm start