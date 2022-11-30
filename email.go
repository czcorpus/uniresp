// Copyright 2022 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2022 Martin Zimandl <martin.zimandl@gmail.com>
// Copyright 2022 Institute of the Czech National Corpus,
//                Faculty of Arts, Charles University
//   This file is part of uniresp.
//
//  uniresp is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  uniresp is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//
//  You should have received a copy of the GNU General Public License
//  along with uniresp.  If not, see <https://www.gnu.org/licenses/>.

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
