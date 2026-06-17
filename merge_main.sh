msg=$1 #"Bug fixes and updates."
type="Fix"
branch="dev"

if [[ -z "${msg}" ]]
then
	msg="Bug fixes and updates."
fi


OPTSTRING="t:m:b:"

while getopts ${OPTSTRING} opt
do
	case ${opt} in
  	t)
    	type=$OPTARG
      	;;
	m)
    	msg=$OPTARG
      	;;
	b)
      	branch=$OPTARG
      	;;
    ?)
      echo "Invalid option: -${OPTARG}."
      exit 1
      ;;
  esac
done

echo "${type}: ${msg}"
echo ${branch}

python scripts/update_version.py


./commit.sh -t "${type}" -m "${msg}"  

git switch main
git merge dev -m "${type}: ${msg}"

#git push -u origin main
./commit.sh -t "${type}" -m "${msg}"

git switch dev
