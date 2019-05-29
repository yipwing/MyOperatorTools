

# try:
#     import xml.etree.cElementTree as ET
# except ImportError:
import xml.etree.ElementTree as ET


def main():
    tree = ET.parse("./sshPort/debug.xml")
    print(tree.getroot())


if __name__ == "__main__":
    main()