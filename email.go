// Copyright 2022 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2022 Martin Zimandl <martin.zimandl@gmail.com>
// Copyright 2022 Institute of the Czech National Corpus,
//                Faculty of Arts, Charles University
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package uniresp

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
)

// DialSmtpServer dials an SMTP server and returns a configured client.
// It supports both insecure and TLS access. In case the port in
// the "server" value is 25 (e.g. "localhost:25") then the insecure
// variant is used. For other port values, the function tries to handle
// TLS connection.
func DialSmtpServer(server, username, password string) (*smtp.Client, error) {
	host, port, err := net.SplitHostPort(server)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SMTP server info: %w", err)
	}
	if port == "25" {
		ans, err := smtp.Dial(server)
		if err != nil {
			return nil, fmt.Errorf("failed to dial: %w", err)
		}
		return ans, err
	}
	auth := smtp.PlainAuth("", username, password, host)
	client, err := smtp.Dial(server)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}
	client.StartTLS(&tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to StartTLS: %w", err)
	}
	err = client.Auth(auth)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}
	return client, nil
}
