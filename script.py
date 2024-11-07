from hugchat import hugchat # type: ignore
from hugchat.login import Login # type: ignore
from youtube_transcript_api import YouTubeTranscriptApi
from googleapiclient.discovery import build
import os
import sys
from dotenv import load_dotenv
load_dotenv()

yt_api_key = os.environ.get("YT_API_KEY")  # Replace 'YT_API_KEY' with your actual environment variable name
huggingface_user = os.environ.get("HUGGINGFACE_USER")  # Replace 'HUGGINGFACE_USER' with the correct name
huggingface_pwd = os.environ.get("HUGGINGFACE_PWD")
video_url = sys.argv[1]
video_id = video_url.split("v=")[1]

youtube = build('youtube','v3',developerKey=yt_api_key)
captions = youtube.captions().list(part='snippet',videoId=video_id).execute()
caption = captions['items'][0]['id']
transcript_list = YouTubeTranscriptApi.get_transcript(video_id)
if not all([yt_api_key, huggingface_user, huggingface_pwd]):
    print("Missing environment variables. Please set 'YT_API_KEY', 'HUGGINGFACE_USER', and 'HUGGINGFACE_PWD'.")
    sys.exit(1)

transcript_txt= ""

for trans in transcript_list :
    transcript_txt += trans['text']


if transcript_txt != "":
      sign = Login(huggingface_user,huggingface_pwd)

      cookies = sign.login()

      cookie_path_dir = "./cookies_snapshot"

      sign.saveCookiesToDir(cookie_path_dir)

      chatbot = hugchat.ChatBot(cookies=cookies.get_dict())
      query_res = chatbot.chat("Summarize in 10 lines: "+transcript_txt) 
      print("Summary:")
      print(query_res)

else:
     print("No Transcript found")         
