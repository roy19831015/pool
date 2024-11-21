DATE=$(date +%Y%m%d%H%M%S)
echo "git commit -a -m '第v1.0.$DATE版本'"
git commit -a -m "第v1.0.$DATE版本"
echo "git checkout -b release/v1.0.$DATE"
git checkout -b release/v1.0.$DATE
echo "git push -u origin release/v1.0.$DATE"
git push -u origin release/v1.0.$DATE
echo "git tag v1.0.$DATE"
git tag v1.0.$DATE
echo "git push --tags"
git push --tags
git checkout master
echo "git checkout master"
