    # Strong Fellas Moving — Official Website

    Professional moving labor assistance website based in Calgary & Surrounding Areas. Features a dynamic cost estimator, AJAX-based lead submission, integrated Telegram notifications, and built-in security rate limiting.

    ## 🛠️ Tech Stack

    - **Backend:** Go (Golang) 1.26.3+
    - **Frontend:** HTML5, Tailwind CSS (via Browser CDN v4), Vanilla JavaScript (AJAX / Fetch API)
    - **Database:** PostgreSQL
    - **Security:** `golang.org/x/time/rate` (Token Bucket Rate Limiter)
    - **Notifications:** Telegram Bot API

    ## 🚀 Features

    1. **Smart Quote Calculator:** Dynamic price calculation based on the number of movers, hours, and addresses.
    2. **AJAX Form Submission:** Smooth user experience without page reloads.
    3. **Telegram Lead Cards:** Instant notifications sent directly to the business owner's Telegram app when a new quote is submitted.
    4. **IP-Based Rate Limiting:** Advanced protection against spam, DDoS, and brute-force attacks by stripping dynamic browser ports (`net.SplitHostPort`).
    5. **Fully Responsive Legal Pages:** Custom designed, independent `Privacy Policy` and `Terms & Conditions` pages.

    ## 💻 Local Setup & Installation

    ### 1. Clone the repository
    ```bash
    git clone [https://github.com/YOUR_USERNAME/strong-fellas-moving.git](https://github.com/YOUR_USERNAME/strong-fellas-moving.git)
    cd strong-fellas-moving