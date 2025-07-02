BASE_URL="http://localhost:17000/?cmd="


curl -s "${BASE_URL}white"
curl -s "${BASE_URL}border%20green"
curl -s "${BASE_URL}figure%20400%20400"
curl -s "${BASE_URL}update"

echo "Commands sent: white background, green border, figure at (400,400)"
