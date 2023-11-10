#!/bin/bash
systemctl daemon-reload
systemctl restart host-api-service
systemctl enable host-api-service
