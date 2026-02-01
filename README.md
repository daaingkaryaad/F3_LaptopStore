# Laptop Store Project

## Team Members
- Ingkar Adilbek
- Syndaly Yerzhan
- Symbat Saparbay

## Overview

This project is part of Assignment 4 for Advanced Programming 1, where we implement the backend for an e-commerce platform specializing in laptops.

## Features

- Detailed Product Descriptions: Laptops with specifications, features, pricing, and availability.
- Advanced Filtering System: Search and filter laptops based on criteria.
- Cart and Order Management: Add items to cart, create orders, and track purchases.

## Requirements

- Backend: Implement an HTTP server with at least 3 working endpoints.
- Data Model: Use Go structs to represent the product catalog, cart, and orders.
- Concurrency: Implement a goroutine for background tasks.
- Git Workflow: Feature branches with at least 2 commits per team member.


### Install dependencies:
go mod tidy

### Running the Backend

## Start the server:
go run .

It will run on http://localhost:8080.

### API Endpoints
GET /api/laptops: Fetch the list of laptops.
POST /api/laptops: Add a new laptop (admin only).
POST /api/cart/items: Add an item to the cart.
GET /api/cart: View the cart.
POST /api/orders: Create an order.

### Demo & Explanation
Demonstrate the working backend and API usage.
Show how data models and features follow the ERD from Assignment 3.

## Setup

### Prerequisites

- Install [Go](https://golang.org/doc/install).

### Installation

Clone the repository:

```bash
git clone https://github.com/daaingkaryaad/F3_LaptopStore.git
cd F3_LaptopStore