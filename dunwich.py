#!/usr/bin/env python3

from urllib.request import urlopen
from http.server import BaseHTTPRequestHandler, HTTPServer
from urllib.parse import urlparse, parse_qs
from multiprocessing import Process, Pool, cpu_count
from random import randint
import os.path, re, copy, html, pprint, time

PORT = 8000
BOOK_LOCAL = "dunwich-horror.html"
BOOK_REMOTE = "https://www.gutenberg.org/files/50133/50133-h/50133-h.htm"

class Book():
    def __init__(self):
        if not os.path.isfile(BOOK_LOCAL):
            print("Grabbing HTML from Gutenberg")
            self.raw_book = urlopen(BOOK_REMOTE).read()
            book_file = open(BOOK_LOCAL, 'wb')
            book_file.write(self.raw_book)
            book_file.close()

        print("Loading HTML")
        book_file = open(BOOK_LOCAL, 'r')
        self.raw_book = book_file.read()
        book_file.close()
        self.__build_chunks__()

    def __build_chunks__(self):
        tmp_book = re.sub(r'<!DOCTYPE[\s\S]*</head>', '', self.raw_book)
        matches = re.findall(r'<p>([\s\S]{30,}?)</p>',tmp_book)

        for index, item in enumerate(matches):
            matches[index] = re.sub(r'<.*?>',"",matches[index])
            matches[index] = matches[index].strip()

        self.chunks = matches
        self.num_chunks = len(matches)
        print("Finished building chunks")

    def getChunk(self, num=0):
        if num < 0 or num >= self.num_chunks:
            return 'Invalid Chunk Index'
        else:
            return self.chunks[num]


class ChunkServer(HTTPServer):
    def __init__(self, *args, book=None, **kw):
        super(HTTPServer,self).__init__(*args,**kw)
        self.num_chunks = book.num_chunks
        self.book = book
        self.requets_handled = 0
        print("Initialized server with "+str(self.num_chunks)+" chunks")


class ChunkHandler(BaseHTTPRequestHandler):

    def do_GET(self):
        target_chunk = int(parse_qs(urlparse(self.path).query).get('chunk')[0])

        self.send_response(200)
        self.send_header('Content-type','text/html')
        self.end_headers()
        # Send the html message
        if target_chunk:
            self.wfile.write(str.encode(self.server.book.getChunk(target_chunk)))
        else:
            self.wfile.write(str.encode("Please specify chunk!"))

        self.server.requets_handled += 1
        return

    def log_request(self, format, *args):
        print("Handled: "+str(self.server.requets_handled))
        return


def client_loop(Book, remote_url):
    result = True
    count = 0

    while result:
        result = getRandomChunk(Book, remote_url)
        count += 1

    return count


def init_as_server(Book):
    server = ChunkServer(("", PORT), ChunkHandler, book=Book)
    print('Started server on ', PORT)
    server.serve_forever()


def init_as_client(Book, remote_url):
    print("Initialized client with "+str(Book.num_chunks)+" chunks")
    pool_size = int(cpu_count() * 1.5)

    process_handlers = []

    for x in range(pool_size):
        p = Process(target=client_loop, args=(Book, remote_url))
        p.start()
        process_handlers.append(p)

    for x in range(pool_size):
        process_handlers[x].join()

def getRandomChunk(Book, remote_url):
    chunk = randint(1,Book.num_chunks-1)
    response = urlopen("http://"+str(remote_url)+":"+str(PORT)+"/?chunk="+str(chunk)).read().decode()

    if response == Book.getChunk(chunk):
        return True
    else:
        print("Retrieved: "+str(chunk))
        print(response)
        print("Have: "+str(chunk))
        print(Book.getChunk(chunk))
        return False

def get_operating_mode():
    print("Please select operating mode:")
    print("1 - Server")
    print("2 - Client")
    print("Selection: ")
    selection = input()

    if selection == '1':
        return 'server'
    elif selection == '2':
        return 'client'
    else:
        return get_operating_mode()

def getIP():
    print('Please enter target IP:')
    ip = input().strip()
    return ip


def main():
    dunny = Book()
    num_chunks = dunny.num_chunks

    mode = get_operating_mode()

    if mode == 'server':
        init_as_server(dunny)
    else:
        ip = getIP()
        init_as_client(dunny, ip)


if __name__ == '__main__':
    main()
