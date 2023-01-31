// This file is autogenerated from scripts/make_python_stdlib.py
// To update this file, update the python versions list in
// scripts/make_python_stdlib.py then run `make python-stdlib`
package frame

var (
	pythonStdlib = map[string]struct{}{
		"__future__":       {},
		"__hello__":        {},
		"__phello__":       {},
		"__phello__.foo":   {},
		"_aix_support":     {},
		"_ast":             {},
		"_bootlocale":      {},
		"_bootsubprocess":  {},
		"_collections_abc": {},
		"_compat_pickle":   {},
		"_compression":     {},
		"_dummy_thread":    {},
		"_markupbase":      {},
		"_osx_support":     {},
		"_py_abc":          {},
		"_pydecimal":       {},
		"_pyio":            {},
		"_sitebuiltins":    {},
		"_strptime":        {},
		"_thread":          {},
		"_threading_local": {},
		"_weakrefset":      {},
		"abc":              {},
		"aifc":             {},
		"antigravity":      {},
		"argparse":         {},
		"array":            {},
		"ast":              {},
		"asynchat":         {},
		"asyncio":          {},
		"asyncore":         {},
		"atexit":           {},
		"audioop":          {},
		"base64":           {},
		"bdb":              {},
		"binascii":         {},
		"binhex":           {},
		"bisect":           {},
		"builtins":         {},
		"bz2":              {},
		"cProfile":         {},
		"calendar":         {},
		"cgi":              {},
		"cgitb":            {},
		"chunk":            {},
		"cmath":            {},
		"cmd":              {},
		"code":             {},
		"codecs":           {},
		"codeop":           {},
		"collections":      {},
		"colorsys":         {},
		"compileall":       {},
		"concurrent":       {},
		"configparser":     {},
		"contextlib":       {},
		"contextvars":      {},
		"copy":             {},
		"copyreg":          {},
		"crypt":            {},
		"csv":              {},
		"ctypes":           {},
		"curses":           {},
		"dataclasses":      {},
		"datetime":         {},
		"dbm":              {},
		"decimal":          {},
		"difflib":          {},
		"dis":              {},
		"distutils":        {},
		"doctest":          {},
		"dummy_threading":  {},
		"email":            {},
		"encodings":        {},
		"ensurepip":        {},
		"enum":             {},
		"errno":            {},
		"faulthandler":     {},
		"fcntl":            {},
		"filecmp":          {},
		"fileinput":        {},
		"fnmatch":          {},
		"formatter":        {},
		"fpectl":           {},
		"fractions":        {},
		"ftplib":           {},
		"functools":        {},
		"gc":               {},
		"genericpath":      {},
		"getopt":           {},
		"getpass":          {},
		"gettext":          {},
		"glob":             {},
		"graphlib":         {},
		"grp":              {},
		"gzip":             {},
		"hashlib":          {},
		"heapq":            {},
		"hmac":             {},
		"html":             {},
		"http":             {},
		"idlelib":          {},
		"imaplib":          {},
		"imghdr":           {},
		"imp":              {},
		"importlib":        {},
		"inspect":          {},
		"io":               {},
		"ipaddress":        {},
		"itertools":        {},
		"json":             {},
		"keyword":          {},
		"lib2to3":          {},
		"linecache":        {},
		"locale":           {},
		"logging":          {},
		"lzma":             {},
		"macpath":          {},
		"macurl2path":      {},
		"mailbox":          {},
		"mailcap":          {},
		"marshal":          {},
		"math":             {},
		"mimetypes":        {},
		"mmap":             {},
		"modulefinder":     {},
		"msilib":           {},
		"msvcrt":           {},
		"multiprocessing":  {},
		"netrc":            {},
		"nis":              {},
		"nntplib":          {},
		"ntpath":           {},
		"nturl2path":       {},
		"numbers":          {},
		"opcode":           {},
		"operator":         {},
		"optparse":         {},
		"os":               {},
		"os2emxpath":       {},
		"ossaudiodev":      {},
		"parser":           {},
		"pathlib":          {},
		"pdb":              {},
		"pickle":           {},
		"pickletools":      {},
		"pipes":            {},
		"pkgutil":          {},
		"platform":         {},
		"plistlib":         {},
		"poplib":           {},
		"posix":            {},
		"posixpath":        {},
		"pprint":           {},
		"profile":          {},
		"pstats":           {},
		"pty":              {},
		"pwd":              {},
		"py_compile":       {},
		"pyclbr":           {},
		"pydoc":            {},
		"pydoc_data":       {},
		"queue":            {},
		"quopri":           {},
		"random":           {},
		"re":               {},
		"readline":         {},
		"reprlib":          {},
		"resource":         {},
		"rlcompleter":      {},
		"runpy":            {},
		"sched":            {},
		"secrets":          {},
		"select":           {},
		"selectors":        {},
		"shelve":           {},
		"shlex":            {},
		"shutil":           {},
		"signal":           {},
		"site":             {},
		"smtpd":            {},
		"smtplib":          {},
		"sndhdr":           {},
		"socket":           {},
		"socketserver":     {},
		"spwd":             {},
		"sqlite3":          {},
		"sre":              {},
		"sre_compile":      {},
		"sre_constants":    {},
		"sre_parse":        {},
		"ssl":              {},
		"stat":             {},
		"statistics":       {},
		"string":           {},
		"stringprep":       {},
		"struct":           {},
		"subprocess":       {},
		"sunau":            {},
		"symbol":           {},
		"symtable":         {},
		"sys":              {},
		"sysconfig":        {},
		"syslog":           {},
		"tabnanny":         {},
		"tarfile":          {},
		"telnetlib":        {},
		"tempfile":         {},
		"termios":          {},
		"test":             {},
		"textwrap":         {},
		"this":             {},
		"threading":        {},
		"time":             {},
		"timeit":           {},
		"tkinter":          {},
		"token":            {},
		"tokenize":         {},
		"tomllib":          {},
		"trace":            {},
		"traceback":        {},
		"tracemalloc":      {},
		"tty":              {},
		"turtle":           {},
		"turtledemo":       {},
		"types":            {},
		"typing":           {},
		"unicodedata":      {},
		"unittest":         {},
		"urllib":           {},
		"uu":               {},
		"uuid":             {},
		"venv":             {},
		"warnings":         {},
		"wave":             {},
		"weakref":          {},
		"webbrowser":       {},
		"winreg":           {},
		"winsound":         {},
		"wsgiref":          {},
		"xdrlib":           {},
		"xml":              {},
		"xmlrpc":           {},
		"zipapp":           {},
		"zipfile":          {},
		"zipimport":        {},
		"zlib":             {},
		"zoneinfo":         {},
	}
)