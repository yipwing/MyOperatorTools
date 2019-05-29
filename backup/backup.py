import sys
import os
import re
import shutil

move_command = 'mv `ls %s|grep -v "%s"` %s && ll'

# TODO finish this.
def MovingPath(path, dest, exception):
    if os.path.exists(path):
        value = raw_input("are you sure you want to move %s to %s ?" % (path, dest))
        if value == 'yes' or value == 'y':
            if os.path.exists(dest):
                status = os.system(move_command % (path, exception, dest))
                print(status)
            else:
                os.makedirs(dest, 755)
                status = os.system(move_command % (path, exception, dest))
                print(status)
    print("the source directory doesn't exists\nmake sure your input source path is exists.")
    
if __name__ == "__main__":
    print(os.path.dirname(sys.argv[2]) + "/" + sys.argv[3])
    try:
        MovingPath(sys.argv[1], sys.argv[2], sys.argv[3])
    except:
        print("please input source directory(full path), destination(full path) and except file. e.g: \nbackup.py /usr/local/2020/01/01 /bk/backup4 1 2020-01-10.log")
    pass
    