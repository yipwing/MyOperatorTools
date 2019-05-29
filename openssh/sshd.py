import os
import time

write_text = '''
auth        required      pam_tally2.so deny=3 even_deny_root unlock_time=600
auth        sufficient    pam_unix.so nullok try_first_pass
auth        required      pam_env.so
auth        required      pam_deny.so

account     required      pam_unix.so
account     required      pam_localuser.so
account     required      pam_permit.so

password    sufficient    pam_unix.so sha512 shadow nullok try_first_pass use_authtok
password    required      pam_deny.so

session     required      pam_limits.so
session     required      pam_unix.so
session     required      pam_lastlog.so nowtmp
'''

sshd_config = '''
Port 3722

AuthorizedKeysFile      .ssh/authorized_keys

PasswordAuthentication yes

UsePAM yes

ClientAliveInterval 150
ClientAliveCountMax 2

Subsystem       sftp    /usr/libexec/sftp-server

PermitRootLogin yes

'''

copy_command = "cp sshd /usr/sbin"

def write():
    # value = raw_input("are you sure you want to move overwrite to those files.(y[yes]/n)")
    # if value == 'yes' or value == 'y':
    try:
        with open("/etc/pam.d/sshd", "w") as f:
            f.write(write_text)
    except Exception as e:
        print(e)
        raise e
    try:
        os.system("cp /etc/ssh/sshd_config /etc/ssh/sshd_config_bak")
        with open("/etc/ssh/sshd_config", "w") as f:
            f.write(sshd_config)
    except Exception as e:
        print(e)
        raise e
    os.system("systemctl stop sshd")
    try:
        time.sleep(1)
        os.system("mv /usr/sbin/sshd /usr/sbin/sshd_bak")
        os.system(copy_command)
        os.system("chmod +x /usr/sbin/sshd")
    except Exception as e:
        print(e)
        raise e
    os.system("systemctl start sshd")
    try:
        time.sleep(1)
        os.system("echo TMOUT=300 >> /root/.bashrc")
        os.system("echo TMOUT=300 >> /root/.zshrc")
    except Exception as e:
        print(e)
        raise e
    return "done"


if __name__ == "__main__":
    print write()