package gpt

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

var (
	OPENAI_API_KEY 	string
	prompt 					string
)

func DetermineUnsubscribeStatus(message string) string {
	initGpt()

	c := openai.NewClient(OPENAI_API_KEY)
	ctx := context.Background()

	prompt += message

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 1,
		Messages:    []openai.ChatCompletionMessage{
			{
				Role: "system",
				Content: "You are a helpful assistant.",
			},
			{
				Role: "user",
				Content: prompt,
			},
		},
	}
	res, err := c.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("Completion error: %v\n", err)
		return ""
	}

	return res.Choices[0].Message.Content
}

func initGpt() {
	OPENAI_API_KEY 	= os.Getenv("OPENAI_API_KEY")

	prompt = `I am programmatically attempting to unsubscribe from a particular sender via there email address. Sometimes when attempting to unsubscribe, the response body can be quite confusing. I need to send you the confusing response messages to determine if unsubscribing was successful for the given sender. After making your determine, return a boolean value of "true" if unsubscribing was successful, and return a boolean value of "false" if unsubscribing was not successful.
	
	Here are some examples:

	1) Message body: "User successfully unsubscribed from campaign: JOB_ALERT". Explanation: The text "User successfully unsubscribed" verifies that unsubscribing was successful. Therefore, a boolean value of "true" should be returned.

	2) Message body: "Your subscription had been successfully adjusted." Explanation: The text "Your subscription had been successfully adjusted" verifies that unsubscribing was successful. Therefore, a boolean value of "true" should be returned.

	3) Message body: "". Explanation: This message is empty which suggests that unsubscribing was not successful. Therefore, a boolean value of "false" should be returned.

	4) Message body:
	
		"<!DOCTYPE html>
			<html lang="en">
					<head>
							<meta charset="utf-8">
							<meta http-equiv="X-UA-Compatible" content="IE=edge">
							<meta name="viewport" content="width=device-width, initial-scale=1">
							<title>Subscription Preferences</title>
							
					

							
					<link rel="stylesheet" href="https://d1dfgjtvrwaror.cloudfront.net/app/releases/2024052001/track/prefs.css">

					</head>
					<body>
							
					

							<div class="container">
									
					<h1>
							<a target="_blank" href="https://www.levels.fyi">
							
									<span class="hidden">Levels.fyi</span>
									<img alt="Levels.fyi" class="logo" src="https://dofdemfeqkgc0.cloudfront.net/account/c48fad78-de6e-4022-8914-f656c4b57d75/account_image/274c0155-3ad7-4f6c-907f-bc3dea2e58a4">
							
							</a>
					</h1>
					
							<p>Your subscription preferences have been updated.</p>
					
							<p>
									<a class="btn" href="https://www.levels.fyi">Return to our website</a>
							</p>
					
					
					

							</div>
							
					</body>
			</html>
		"

	Explanation: This message may appear confusing initially, but upon closer inspection, we can find the text "Your subscription preferences have been updated." Therefore, a boolean value of "true" should be returned.

	5) Message body:
		<!DOCTYPE html>

				<link rel="shortcut icon" href="/media/Favicon-16by16.png"/>

			<html lang="en" style="height: 100%">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Consent Pages</title>
				<script nonce="bm1wig7kuX4rdLnSDJIaDA==">
					function handleSubmit(e) {
						e.preventDefault();
						// using XMLHttpRequest for browser compatibility
						var xhr = new XMLHttpRequest();
						xhr.open("POST", window.location.href, true);
						xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
						xhr.onreadystatechange = function () { // Call a function when the state changes.
							if (xhr.readyState === XMLHttpRequest.DONE) {
								if (xhr.status === 200) {
									// # move to unsubscribe success
									document.getElementById('unsubscribe-fallback').style.display = 'none';
									document.getElementById('unsubscribe-success-fallback').style.display = 'block';
								} else {
									alert('Something went wrong. Please try again.')
								}
							}
						}
						xhr.send(new URLSearchParams(new FormData(document.getElementById('unsubscribe-form'))).toString());
					}
				</script>
			</head>
			<main style="height: 100%;">
				<body style="margin:0; height: 100%;">
					<div id="onsite-consent-pages-root" style="height: 100%"></div>
					<script nonce="bm1wig7kuX4rdLnSDJIaDA==" id="onsite-js-bundle" type="text/javascript" src="https://static.kmail-lists.com/onsite-consent-pages/js/onsiteConsentPages.js"></script>
					<script nonce="bm1wig7kuX4rdLnSDJIaDA==">
						document.addEventListener('DOMContentLoaded', function () {
							document.getElementById('onsite-js-bundle')
								.addEventListener('error', function () {
									window.location.search.indexOf('failedJS') === -1 ? window.location.href += '&failedJS=true' : null
								});
						});
					</script>
					<div id="unsubscribe-fallback" style="height: 100%; display: none;">
						<div
							style="background-color: rgb(245, 245, 245); display: flex; justify-content: center; min-height: 100%; padding-top: 40px; padding-bottom: 40px; width: 100%; box-sizing: border-box;">
							<form id="unsubscribe-form" style="max-width: 100%;" method="POST">
								<div
									style="height: fit-content; width: 558px; background-color: rgb(255, 255, 255); border: 1px none rgb(204, 204, 204); border-radius: 0px; padding: 60px; max-width: 100%; box-sizing: border-box;">
									<div style="display: flex; position: relative;">
										<div style="display: flex; flex-direction: column; padding: 0px; position: relative; flex: 1 1 0%;">
											<style>* {
												margin-block: unset;
											}</style>
											<div style="width: 100%; font-family: &quot;Helvetica Neue&quot;, Arial; color: rgb(34, 34, 34);"><h1
												style="text-align: center;">Unsubscribe</h1></div>
										</div>
									</div>
									<div style="display: flex; position: relative;">
										<div
											style="display: flex; flex-direction: column; padding: 10px 6px; background-color: rgba(255, 255, 255, 0); position: relative; flex: 1 1 0%;">
											<div style="display: flex; width: 100%; flex-direction: column;"><label
												for="kl-consent-page-7eafd389433e4c48b01ba3a907d6246c"
												style="font-family: &quot;Helvetica Neue&quot;, Arial; font-size: 16px; font-weight: 700; color: rgb(34, 34, 34); margin-bottom: 8px;">Enter the email you want to unsubscribe<span
												style="color: rgb(231, 76, 60);">*</span></label>
												<div style="display: flex;"><input
													id="kl-consent-page-7eafd389433e4c48b01ba3a907d6246c"
													aria-required="true"
													value=""
													name="$email"
													type="email"
													style="color: rgb(34, 34, 34); font-family: &quot;Helvetica Neue&quot;, Arial; font-size: 16px; font-weight: 400; width: 100%; height: 50px; padding-left: 16px; border: 1px solid rgb(204, 204, 204); border-radius: 0px; box-sizing: border-box;">
												</div>
											</div>
										</div>
									</div>
									<div style="display: flex; position: relative;">
										<div
											style="display: flex; flex-direction: column; padding: 10px 6px; background-color: rgba(255, 255, 255, 0); position: relative; flex: 1 1 0%;">
											<button type="submit"
															style="font-family: &quot;Helvetica Neue&quot;, Arial; font-size: 16px; font-weight: 400; color: rgb(255, 255, 255); background-color: rgb(17, 85, 204); border-radius: 0px; border: 0px none rgb(34, 34, 34); height: 50px; width: 100%; padding: 0px; margin: 0px auto; cursor: pointer;">
												Unsubscribe
											</button>
										</div>
									</div>
									<div style="display: flex; position: relative;">
										<div
											style="display: flex; flex-direction: column; padding: 10px 6px; background-color: rgba(255, 255, 255, 0); position: relative; flex: 1 1 0%;"></div>
									</div>
								</div>
							</form>
						</div>
					</div>
					<div id="unsubscribe-success-fallback" style="height: 100%; display: none;">
						<div style="height: 100%">
							<div
								style="background-color: rgb(245, 245, 245); display: flex; justify-content: center; min-height: 100%; padding-top: 40px; padding-bottom: 40px; width: 100%; box-sizing: border-box;">
								<div
									style="height: fit-content; width: 558px; background-color: rgb(255, 255, 255); border: 1px none rgb(204, 204, 204); border-radius: 0px; padding: 60px; max-width: 100%; box-sizing: border-box;">
									<div style="display: flex; position: relative;">
										<div
											style="display: flex; flex-direction: column; padding: 10px 6px; background-color: rgb(255, 255, 255); position: relative; flex: 1 1 0%;">
											<style>* {
												margin-block: unset;
											}</style>
											<div style="width: 100%; font-family: &quot;Helvetica Neue&quot;, Arial; color: rgb(34, 34, 34);"><h1
												style="text-align: center;">You've been unsubscribed</h1></div>
										</div>
									</div>
								</div>
							</div>
						</div>
					</div>
				</body>
				<script nonce="bm1wig7kuX4rdLnSDJIaDA==" type="application/javascript">
					if (window.location.search.indexOf('failedJS') !== -1
							|| window.location.search.indexOf('failedAPIRequest') !== -1) {
						if (window.location.href.indexOf('subscriptions/unsubscribed') !== -1) {
							document.getElementById('onsite-consent-pages-root').style.display = 'none';
							document.getElementById('unsubscribe-success-fallback').style.display = 'block';
						} else if (window.location.href.indexOf('subscriptions/unsubscribe') !== -1) {
							document.getElementById('unsubscribe-fallback').style.display = 'block';
							document.getElementById('onsite-consent-pages-root').style.display = 'none';
							document.getElementById('unsubscribe-form').onsubmit = handleSubmit;
						}
					}
				</script>
			</main>
			</html>
		"

	Explanation: this message may appear especially confusing initially, but upon further inspection, we can locate the text "You've been unsubscribed". Therefore, a boolean value of "true" should be returned.

	6) Message body: "May need to manually open link". Explanation: This message suggests that further action is required from the user. Therefore, a boolean value of "fakse" should be returned.

	Note 1: the above are merely examples and in reality, the message bodies can vary greatly. Therefore, it's important to carefully inspect each message body as a unique case before making a final determination.

	Note 2: if at all doubtful, it's better to be safe. Therefore, only return a boolean value of "true" if absolutely certain that unsubscribed was successful.

	With the above context given, please make a determination on the following message by only returning either "true" or "false". Do NOT give your explanation. I will reiterate, ONLY return a "true" or a "false": 
	`
}
