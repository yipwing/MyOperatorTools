import os
import sys
import subprocess


ipAddr = [
    "172.17.113.42",
    "172.17.113.41",
    "172.17.113.44",
    "172.17.113.43",
    "172.17.103.227",
    "172.17.103.228",
    "172.17.103.229",
    "172.17.103.231",
    "172.17.103.232",
    "172.17.103.233",
    "172.17.103.234",
    "172.17.103.235",
    "172.17.101.129",
    "172.17.101.128",
    "172.17.107.44",
    "172.17.103.230",
    "172.17.103.237",
    "172.17.101.127",
    "172.17.101.126",
    "172.17.113.38",
    "172.17.113.37",
    "172.17.113.36",
    "172.17.113.35",
    "172.17.113.34",
    "172.17.113.33",
    "172.17.113.32",
    "172.17.113.31",
    "172.17.101.125",
    "172.17.101.124",
    "172.17.101.123",
    "172.18.1.222",
    "172.17.107.43",
    "172.17.107.42",
    "172.17.107.41",
    "172.17.107.40",
    "172.17.107.39",
    "172.17.107.38",
    "172.17.107.37",
    "172.17.107.36",
    "172.17.107.35",
    "172.17.107.34",
    "172.17.113.30",
    "172.17.113.29",
    "172.17.113.27",
    "172.17.113.28",
    "172.17.113.25",
    "172.17.113.24",
    "172.17.113.22",
    "172.17.113.21",
    "172.17.103.236",
    "172.17.107.33",
    "172.17.107.32",
    "172.17.107.31",
    "172.17.107.30",
    "172.17.107.29",
    "172.17.107.27",
    "172.17.107.26",
    "172.17.107.25",
    "172.17.107.24",
    "172.17.103.238",
    "172.17.103.239",
    "172.17.100.117",
    "172.17.100.116",
    "172.17.101.116",
    "172.17.101.115",
    "172.17.100.115",
    "172.17.100.114",
    "172.17.103.240",
    "172.17.105.23",
    "172.17.105.25",
    "172.17.103.241",
    "172.17.103.242",
    "172.17.103.243",
    "172.17.107.23",
    "172.17.107.22",
    "172.17.105.26",
    "172.17.107.21",
    "172.17.107.20",
    "172.17.105.27",
    "172.17.100.93",
    "172.17.101.87",
    "172.17.107.19",
    "172.17.103.245",
    "172.17.107.13",
    "172.17.107.12",
    "172.17.107.11",
    "172.17.107.10",
    "172.17.107.15",
    "172.17.107.14",
    "172.17.107.8",
    "172.17.107.7",
    "172.17.107.17",
    "172.17.107.18",
    "172.17.107.4",
    "172.17.107.3",
    "172.17.103.246",
    "172.17.103.247",
    "172.17.103.248",
    "172.17.105.20",
    "172.17.103.249",
    "172.17.105.18",
    "172.17.105.17",
    "172.17.105.13",
    "172.17.101.122",
    "172.17.101.121",
    "172.17.101.120",
    "172.17.101.112",
    "172.17.101.110",
    "172.17.101.111",
    "172.17.101.109",
    "172.17.113.20",
    "172.17.113.19",
    "172.17.113.18",
    "172.17.113.17",
    "172.17.113.16",
    "172.17.113.14",
    "172.17.113.15",
    "172.17.113.13",
    "172.17.113.12",
    "172.17.113.11",
    "172.17.113.10",
    "172.17.113.9",
    "172.17.113.8",
    "172.17.113.7",
    "172.17.113.6",
    "172.17.113.5",
    "172.17.100.108",
    "172.17.101.105",
    "172.17.101.104",
    "172.17.100.103",
    "172.17.102.30",
    "172.17.100.105",
    "172.17.100.106",
    "172.17.100.104",
    "172.17.100.101",
    "172.17.100.102",
    "172.18.1.220",
    "172.17.100.113",
    "172.17.100.112",
    "172.17.101.92",
    "172.17.100.94",
    "172.17.102.13",
    "172.17.102.12",
    "172.17.102.11",
    "172.17.102.9",
    "172.17.102.48",
    "172.17.101.48",
    "172.17.100.48",
    "172.17.102.46",
    "172.17.101.86",
    "172.17.100.49",
    "172.17.102.49",
    "172.17.102.41",
    "172.17.102.40",
    "172.17.101.5",
    "172.17.101.233",
    "172.17.100.89",
    "172.17.100.99",
    "172.17.100.98",
    "172.17.100.96",
    "172.17.100.95",
    "172.17.100.94",
    "172.17.100.91",
    "172.17.100.90",
    "172.17.100.88",
    "172.17.100.87",
    "172.17.101.85",
    "172.17.101.84",
    "172.17.101.83",
    "172.17.101.82",
    "172.17.101.81",
    "172.17.101.80",
    "172.17.101.79",
    "172.17.101.78",
    "172.17.101.77",
    "172.17.101.76",
    "172.17.101.75",
    "172.17.101.74",
    "172.17.101.73",
    "172.17.101.72",
    "172.17.101.71",
    "172.17.101.70",
    "172.17.101.69",
    "172.17.101.68",
    "172.17.101.67",
    "172.17.101.66",
    "172.17.101.64",
    "172.17.101.63",
    "172.17.101.62",
    "172.17.101.61",
    "172.17.101.60",
    "172.17.101.59",
    "172.17.101.58",
    "172.17.101.57",
    "172.17.101.55",
    "172.17.101.54",
    "172.17.100.111",
    "172.17.100.110",
    "172.17.101.51",
    "172.17.101.50",
    "172.17.101.49",
    "172.17.102.29",
    "172.17.102.8",
    "172.17.102.7",
    "172.17.102.6",
    "172.17.102.5",
    "172.17.102.4",
    "172.17.102.3",
    "172.17.102.2",
    "172.17.105.29",
    "172.17.103.226",
    "172.17.103.225",
    "172.17.103.224",
    "172.17.103.223",
    "172.17.103.222",
    "172.17.105.30",
    "172.17.105.31",
    "172.17.105.32",
]

Passwords = [
    "fanjie",
    "Rootadmin"
]

scp_command = "scp -P 3722 ~/ssh_config/* root@%s:/root"
password_enter = ""
excute_command = "python /root/sshd.py" 
ssh_login = "root@%s"

def Batfile():
    for ip in ipAddr:
        try:
            os.system(scp_command % ip)
        except Exception as e:
            print(e)
        ssh_cli = subprocess.Popen(["ssh", "-p 3722", ssh_login % ip], shell=False, stdout=subprocess.PIPE, stderr=subprocess.PIPE, stdin=subprocess.PIPE)
        result = ssh_cli.stdout.readlines()
        if result == []:
            error = ssh_cli.stderr.readlines()
            print >>sys.stderr, "ERROR: %s" % error
        else:
            print (result)
    print("done!")

if __name__ == "__main__":
    Batfile()