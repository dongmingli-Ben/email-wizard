from typing import List, Dict
from email import policy
from email.parser import BytesParser
from email.message import EmailMessage
from bs4 import BeautifulSoup


def extract_plain_text_from_html(html_str: str) -> str:
    soup = BeautifulSoup(html_str, "html.parser")
    return soup.get_text(separator=" ", strip=True)


def is_html(text: str) -> bool:
    soup = BeautifulSoup(text, "html.parser")
    return len(soup.find_all()) > 0


def extract_plain_text(raw_texts: List[str]) -> str:
    texts = []
    for raw_text in raw_texts:
        if is_html(raw_text):
            texts.append(extract_plain_text_from_html(raw_text))
        else:
            texts.append(raw_text)
    return " ".join(texts)


def clean_special_characters(text: str) -> str:
    special_chars = ["\r", "\n", "\xa0"]
    for c in special_chars:
        text = text.replace(c, " ")
    subtexts = text.split()
    subtexts = list(map(lambda s: s.strip(), subtexts))
    subtexts = list(filter(len, subtexts))
    return " ".join(subtexts)


def clean_email_from_raw_texts(raw_email: Dict[str, str]) -> dict:
    email_info = raw_email.copy()

    content = extract_plain_text(raw_email["content"])
    text = clean_special_characters(content)
    email_info["content"] = text

    return email_info
