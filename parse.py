from pdfminer.high_level import extract_text
from sklearn.feature_extraction.text import CountVectorizer
from sklearn.metrics.pairwise import cosine_similarity
import sys

resume = extract_text(sys.argv[1])
desc= sys.argv[2]
vectorize = CountVectorizer().fit_transform([resume,desc])
cosine_sim = cosine_similarity(vectorize)
print(round(cosine_sim[0][1]*100,2))

