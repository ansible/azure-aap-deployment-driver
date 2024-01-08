#!/usr/bin/env python3

import argparse
import math
import requests
import time
import urllib.parse

# Removes all SSO clients from selected server by org ID older than xxx


DEFAULT_ORG = 123
DEFAULT_HOST = "sso.redhat.com"
DEFAULT_URL = f"https://{DEFAULT_HOST}"
TOKEN_PATH = "auth/realms/redhat-external/protocol/openid-connect/token"
API_PATH = "auth/realms/redhat-external/apis/beta/acs/v1"
CLIENT_ID = "aoc-client-manager"


def parse_args():
    parser = argparse.ArgumentParser(description="Remove extraneous SSO dynamic clients older than 4 hours.")

    parser.add_argument(
        "-s",
        "--secret",
        help="Client secret",
        required=True,
    )
    return parser.parse_args()


def load_clients(token):
    url = urllib.parse.urljoin(DEFAULT_URL, API_PATH)
    headers = {}
    headers["Authorization"] = f"Bearer {token}"
    params = {}
    params["orgId"] = DEFAULT_ORG
    resp = requests.get(url, headers=headers, params=params)
    resp.raise_for_status()
    return resp.json()


def clean_clients(token, clients):
    base_url = urllib.parse.urljoin(DEFAULT_URL, API_PATH)
    headers = {}
    headers["Authorization"] = f"Bearer {token}"
    cutoff_time = math.ceil(time.time()) - 14400  # Now - four hours
    for client in clients:
        client_id = client["clientId"]
        if client["createdAt"] < cutoff_time:
            client_url = f"{base_url}/{client_id}"
            print(f"Deleting client ID {client_id}")
            resp = requests.delete(client_url, headers=headers)
            resp.raise_for_status()
        else:
            print(f"Not deleting client ID {client_id} since it is less than 4 hours old")


def get_oauth_token(secret):
    url = urllib.parse.urljoin(DEFAULT_URL, TOKEN_PATH)
    headers = {}
    headers["Host"] = DEFAULT_HOST
    headers["Content-Type"] = "application/x-www-form-urlencoded"
    data = {}
    data["grant_type"] = "client_credentials"
    data["scope"] = "api.iam.clients.aoc"
    data["client_id"] = CLIENT_ID
    data["client_secret"] = secret
    resp = requests.post(url, headers=headers, data=data)
    resp.raise_for_status()
    body = resp.json()
    return body["access_token"]


if __name__ == "__main__":
    args = parse_args()
    token = get_oauth_token(args.secret)
    clients = load_clients(token)
    if not clients:
        print(f"No clients currently exist in org {DEFAULT_ORG}")
    clean_clients(token, clients)
