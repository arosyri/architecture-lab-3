BASE_URL="http://localhost:17000/"

curl -X POST -d "green" $BASE_URL
curl -X POST -d "figure 400 400" $BASE_URL
curl -X POST -d "update" $BASE_URL

x=100
y=100
dx=20
dy=15

while true; do
  curl -X POST -d "move $x $y" $BASE_URL
  curl -X POST -d "update" $BASE_URL
  echo "Moved to ($x,$y)"

  x=$((x + dx))
  y=$((y + dy))

  if [ $x -lt 0 ] || [ $x -gt 800 ]; then
    dx=$(( -dx ))
  fi

  if [ $y -lt 0 ] || [ $y -gt 800 ]; then
    dy=$(( -dy ))
  fi

  sleep 1
done
