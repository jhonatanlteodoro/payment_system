'''
Hello there!
Welcome to data load hell :)

Basically I need to load 4 files, 2 of them are pretty safe to consume
but there's two that is petty big posts[~100GB], and votes [~23GB].

this is script should work fine but will take some time to load the posts because its too big,
off course, assuming you are loading this to a pgsql in a docker file with minimal setup as i am doing here.

Feel free to change the amount of workers, connections and everything, but keep in mind that some queries may fail if youre
database starts to suffer :)

In case you want more data, Link for public files:
https://archive.org/download/stackexchange
'''
import datetime
import random
import time
import xml.etree.ElementTree as ET

import queue as pQueue
from psycopg_pool import ConnectionPool
import multiprocessing

default_date = datetime.datetime.strptime("2000-01-01T00:00:00.000","%Y-%m-%dT%H:%M:%S.%f")
def prepare_string(text):
    if text is None or len(text) == 0:
        return "NOT-DEFINED"
    return escape_sql_string(text)

def escape_sql_string(text):
    return f"{text.replace("'", "''")}"

# Check the URL
db_url = "postgresql://secret_user:secret_password@localhost:5432/stackoverflow"
connection_pool = ConnectionPool(db_url, max_size=5, open=True, num_workers=5)
connection_pool.wait()
print("pool ready")

def main():
    task_queue = multiprocessing.Queue()
    workers = []
    for i in range(10):
        worker = multiprocessing.Process(target=execute_query, args=(task_queue, i))
        workers.append(worker)
        worker.start()
        print(f"started worker {i}")

    populate_users(task_queue, "Users.xml")
    populate_votes(task_queue, "Votes.xml")
    populate_tags(task_queue, "Tags.xml")
    populate_users(task_queue, "Users.xml")
    populate_posts(task_queue, "Posts.xml")
    print("Waiting for workers to finish...")
    for worker in workers:
        worker.join()

    print("All done!")



def execute_query(queue, worker_id):
    while True:
        try:
            query = queue.get(timeout=5)
            with connection_pool.connection() as conn:
                conn.autocommit = True
                # for some reason if I did not set autocommit as true some errors are raised
                # because of uncommited query
                with conn.cursor() as curs:
                    curs.execute(query)
                    print(f"worker {worker_id} made a query")
                time.sleep(random.randint(1,5))
        except pQueue.Empty:
            print(f"worker {worker_id} done")
            break

        except Exception as error:
            print(f"worker {worker_id} failed {error}")
            # it should happen when the database are suffering to process the data
            # from my tests

def populate_users(task_queue, filename):
    for batchNum, batch in enumerate(parse_with_generator(filename)):
        query = "INSERT INTO users (id, display_name, location, reputation, views, up_votes, down_votes) VALUES "
        for idx, item in enumerate(batch):
            query += f"({int(item.get("Id", 1))}, '{prepare_string(item.get("DisplayName"))}', '{prepare_string(item.get("Location"))}', {int(item.get("Reputation", 0))}, {int(item.get("Views", 0))}, {int(item.get("UpVotes", 0))}, {int(item.get("DownVotes", 0))})"
            if idx < len(batch) - 1:
                query += ", "
                continue

        task_queue.put(query)
        print(f"{batchNum+1} Batch users send to worker")

def populate_votes(task_queue, filename):
    for batchNum, batch in enumerate(parse_with_generator(filename)):
        query = "INSERT INTO votes (id, post_id, vote_type_id, creation_date) VALUES "
        for idx, item in enumerate(batch):
            query += f"({int(item.get("Id", 1))}, {int(item.get("PostId", 0))}, {int(item.get("VoteTypeId", 0))}, '{item.get("CreationDate", default_date)}')"
            if idx < len(batch) - 1:
                query += ", "
                continue

        task_queue.put(query)
        print(f"{batchNum + 1} Batch votes send to worker")

def populate_tags(task_queue, filename):
    for batchNum, batch in enumerate(parse_with_generator(filename)):
        query = "INSERT INTO tags (id, name, count, excerpt_post_id, wiki_post_id) VALUES "
        for idx, item in enumerate(batch):
            query += f"({int(item.get("Id", 1))}, '{item.get("TagName")}', {int(item.get("Count", 0))}, {int(item.get("ExcerptPostId", 0))}, {int(item.get("WikiPostId", 0))})"
            if idx < len(batch) - 1:
                query += ", "
                continue

        task_queue.put(query)
        print(f"{batchNum + 1} Batch tags send to worker")

def populate_posts(task_queue, filename):
    for batchNum, batch in enumerate(parse_with_generator(filename, batch_size=1000)):
        query = "INSERT INTO posts (id, post_type_id, creation_date, score, view_count, owner_user_id, tags, answer_count, comment_count, favorite_count) VALUES "
        for idx, item in enumerate(batch):
            query += (f""
                      f"({int(item.get("Id", 1))}, '{item.get("PostTypeId")}', '{item.get("CreationDate", default_date)}',"
                      f"{int(item.get("Score", 0))}, {int(item.get("ViewCount", 0))}, {int(item.get("OwnerUserId", 0))},"
                      f"'{prepare_string(item.get("Tags"))}', {int(item.get("AnswerCount", 0))},"
                      f"{int(item.get("CommentCount", 0))}, {int(item.get("FavoriteCount", 0))})")
            if idx < len(batch) - 1:
                query += ", "
                continue

        task_queue.put(query)
        print(f"{batchNum + 1} Batch posts send to worker")

def parse_with_generator(file_path, batch_size=100):
    batch = []
    for event, elem in ET.iterparse(file_path, events=('end',)):
        if elem.tag == 'row':
            if len(batch) < batch_size:
                batch.append(elem.attrib)
                elem.clear()
                continue

            yield batch
            batch = []

    return batch

if __name__ == '__main__':
    main()