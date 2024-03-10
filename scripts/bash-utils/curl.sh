#!/bin/bash
echo off > /dev/null

function curl_post() {

    result=$(curl -k --location --request POST "$CURL_BASE_URL$1" \
             --header "Zixel-Application-Id: $LOGIN_APP_ID" \
             --header "Zixel-Open-Id: $LOGIN_OPEN_ID" \
             --header "Zixel-Auth-Token: $LOGIN_AUTH_TOKEN" \
             --header "Content-Type: application/json" \
             --header "Accept: */*" \
             --header "Connection: keep-alive" \
             --data-raw "$2")

   echo $result
}

function curl_post_no_openid() {

    result=$(curl -k --location --request POST "$CURL_BASE_URL$1" \
             --header "Zixel-Application-Id: $LOGIN_APP_ID" \
             --header "Zixel-Auth-Token: $LOGIN_AUTH_TOKEN" \
             --header "Content-Type: application/json" \
             --header "Accept: */*" \
             --header "Connection: keep-alive" \
             --data-raw "$2")

   echo $result
}


function curl_get() {

    result=$(curl -k --location --request GET "$CURL_BASE_URL$1" \
             --header "Zixel-Application-Id: $LOGIN_APP_ID" \
             --header "Zixel-Open-Id: $LOGIN_OPEN_ID" \
             --header "Zixel-Auth-Token: $LOGIN_AUTH_TOKEN" \
             --header "Content-Type: application/json" \
             --header "Accept: */*" \
             --header "Connection: keep-alive")

   echo $result
}

function curl_delete() {

    result=$(curl -k --location --request DELETE "$CURL_BASE_URL$1" \
             --header "Zixel-Application-Id: $LOGIN_APP_ID" \
             --header "Zixel-Open-Id: $LOGIN_OPEN_ID" \
             --header "Zixel-Auth-Token: $LOGIN_AUTH_TOKEN" \
             --header "Content-Type: application/json" \
             --header "Accept: */*" \
             --header "Connection: keep-alive")

   echo $result
}

function curl_put() {

    result=$(curl -k --location --request PUT "$CURL_BASE_URL$1" \
             --header "Zixel-Application-Id: $LOGIN_APP_ID" \
             --header "Zixel-Open-Id: $LOGIN_OPEN_ID" \
             --header "Zixel-Auth-Token: $LOGIN_AUTH_TOKEN" \
             --header "Content-Type: application/json" \
             --header "Accept: */*" \
             --header "Connection: keep-alive" \
             --data-raw "$2")

   echo $result
}
