"""Generates a bunch of mock users which can be used to test app during dev"""

from playwright.sync_api import sync_playwright

# Randomly generated data
form_data_list = [
    {
        "username": "fbritnell0",
        "display_name": "Francine Britnell",
        "email": "fbritnell0@un.org",
        "password": "lN1%xYR4)ys)~XE",
    },
    {
        "username": "avreiberg1",
        "display_name": "Ami Vreiberg",
        "email": "avreiberg1@amazon.co.uk",
        "password": "kV3#,t4!'i",
    },
    {
        "username": "jvaladez2",
        "display_name": "Jeffy Valadez",
        "email": "jvaladez2@jalbum.net",
        "password": "bB1*7m)u!9fkE",
    },
    {
        "username": "mcroce3",
        "display_name": "Morena Croce",
        "email": "mcroce3@trellian.com",
        "password": "uY7?T52ba",
    },
    {
        "username": "ahirtzmann4",
        "display_name": "Adam Hirtzmann",
        "email": "ahirtzmann4@wikispaces.com",
        "password": "eZ0\\%9%T(N8#o4",
    },
    {
        "username": "acollecott5",
        "display_name": "Arluene Collecott",
        "email": "acollecott5@harvard.edu",
        "password": "lX4$5bw'w",
    },
    {
        "username": "kriby6",
        "display_name": "Kirsti Riby",
        "email": "kriby6@chronoengine.com",
        "password": "kO1?{Va37P",
    },
    {
        "username": "mmarusik7",
        "display_name": "Mireielle Marusik",
        "email": "mmarusik7@ucoz.com",
        "password": "rQ3=RFOgLnR`B~",
    },
    {
        "username": "abolderstone8",
        "display_name": "Archibaldo Bolderstone",
        "email": "abolderstone8@google.com.hk",
        "password": "qB0{Q05g(5",
    },
    {
        "username": "jdurnford9",
        "display_name": "Janie Durnford",
        "email": "jdurnford9@icio.us",
        "password": 'jF6{@cyGk"m4a',
    },
]

with sync_playwright() as p:
    browser = p.firefox.launch(headless=True)

    for i, form_data in enumerate(form_data_list, start=1):
        print(f"Creating user {i}/{len(form_data_list)}")
        
        page = browser.new_page()
        page.goto("http://localhost:5173/register", wait_until="domcontentloaded")

        page.wait_for_selector("#register-form")

        page.fill("#username", form_data["username"])
        page.fill("#display_name", form_data["display_name"])
        page.fill("#email", form_data["email"])
        page.fill("#password", form_data["password"])
        page.click("#submit-btn")

        page.wait_for_timeout(2000)
        page.close()

    print("All users created")
