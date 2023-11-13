#!/bin/bash
systemctl daemon-reload
systemctl restart host-gpu-api-service
systemctl enable host-gpu-api-service
