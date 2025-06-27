package services

import (
	"fmt"
	"time"
)

func GenerateWelcomeHTML(recipientEmail string) string {
	return fmt.Sprintf(`
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Welcome to ImaginAI!</title>
            <style>
                body {
                    font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
                    line-height: 1.6;
                    color: #333333;
                    background-color: #f7f7f7;
                    margin: 0;
                    padding: 0;
                }
                .container {
                    max-width: 500px;
                    margin: 30px auto;
                    background: #ffffff;
                    border-radius: 8px;
                    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
                    padding: 30px;
                    border: 1px solid #e0e0e0;
                }
                h2 {
                    color: #1a1a1a;
                    font-size: 28px;
                    margin-bottom: 20px;
                    text-align: center;
                }
                p {
                    margin-bottom: 15px;
                }
                .button-container {
                    text-align: center;
                    margin: 30px 0;
                }
                .button {
                    display: inline-block;
                    background-color: #28a745; /* A pleasant green for action */
                    color: #ffffff;
                    padding: 12px 25px;
                    border-radius: 5px;
                    text-decoration: none;
                    font-weight: bold;
                    font-size: 16px;
                    transition: background-color 0.3s ease;
                }
                .button:hover {
                    background-color: #218838; /* Darker green on hover */
                }
                .footer {
                    margin-top: 30px;
                    font-size: 0.9em;
                    color: #777777;
                    text-align: center;
                    border-top: 1px solid #eeeeee;
                    padding-top: 20px;
                }
                a {
                    color: #007bff;
                    text-decoration: none;
                }
                a:hover {
                    text-decoration: underline;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h2>Welcome to ImaginAI!</h2>
                <p>Hello <strong>%s</strong>,</p>
                <p>Thank you for joining the ImaginAI community! We're thrilled to have you on board.</p>
                <p>At ImaginAI, we empower you to unleash your creativity with AI-powered ai chat agents.</p>
                <br/><br/>
                <p>If you have any questions or need assistance, our support team is always here to help. Feel free to reply to this email or visit our help center.</p>
                <p>Happy creating!<br/>The ImaginAI Team</p>
            </div>
            <div class="footer">
                <p>&copy; %d ImaginAI. All rights reserved.</p>
                <p><a href="[Your Website Link]">Our Website</a> | <a href="[Your Privacy Policy Link]">Privacy Policy</a> | <a href="[Your Support Link]">Support</a></p>
            </div>
        </body>
        </html>
    `, recipientEmail, time.Now().Year())
}

func GeneratePasswordResetHTML(code string, recipientEmail string) string {
	return fmt.Sprintf(`
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Password Reset Request</title>
            <style>
                body {
                    font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
                    line-height: 1.6;
                    color: #333333;
                    background-color: #f7f7f7;
                    margin: 0;
                    padding: 0;
                }
                .container {
                    max-width: 500px;
                    margin: 30px auto;
                    background: #ffffff;
                    border-radius: 8px;
                    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
                    padding: 30px;
                    border: 1px solid #e0e0e0;
                }
                h2 {
                    color: #1a1a1a;
                    font-size: 24px;
                    margin-bottom: 20px;
                    text-align: center;
                }
                p {
                    margin-bottom: 15px;
                }
                .code-section {
                    text-align: center;
                    margin: 25px 0;
                }
                .code {
                    font-size: 2.2em;
                    font-weight: bold;
                    letter-spacing: 4px;
                    color: #007bff; /* A nice blue color */
                    background: #e9f5ff; /* Lighter blue background for the code */
                    padding: 15px 30px;
                    border-radius: 8px;
                    display: inline-block;
                    border: 1px dashed #a0d8ff;
                }
                .footer {
                    margin-top: 30px;
                    font-size: 0.9em;
                    color: #777777;
                    text-align: center;
                    border-top: 1px solid #eeeeee;
                    padding-top: 20px;
                }
                .important-note {
                    color: #dc3545; /* Red for important notes */
                    font-weight: bold;
                    margin-top: 20px;
                    padding: 10px;
                    background-color: #ffebeb;
                    border-left: 5px solid #dc3545;
                }
                a {
                    color: #007bff;
                    text-decoration: none;
                }
                a:hover {
                    text-decoration: underline;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h2>Password Reset Request</h2>
                <p>Hello,</p>
                <p>We received a request to reset the password for your account associated with the email address: <strong>%s</strong>.</p>
                <p>To proceed with your password reset, please use the following verification code:</p>
                <div class="code-section">
                    <span class="code">%s</span>
                </div>
                <p>This code is valid for a limited time (e.g., 10 minutes). Please do not share this code with anyone.</p>
                <div class="important-note">
                    If you did not request a password reset, please ignore this email. Your password will remain unchanged.
                </div>
                <p>If you have any questions or encounter issues, please contact our support team.</p>
                <p>Thanks,<br/>The ImaginAI Team</p>
            </div>
            <div class="footer">
                <p>&copy; %d ImaginAI. All rights reserved.</p>
                <p><a href="[Your Website Link]">Our Website</a> | <a href="[Your Privacy Policy Link]">Privacy Policy</a></p>
            </div>
        </body>
        </html>
    `, recipientEmail, code, time.Now().Year())
}
