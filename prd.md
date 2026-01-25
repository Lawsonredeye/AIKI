AIKI – MVP PRODUCT REQUIREMENTS DOCUMENT (PRD)
1. Product Overview
   Product Name: Aiki
   Goal: To help job seekers and young professionals stay productive, consistent, and accountable while searching for jobs or improving their career goals.
   Tagline: “Lock in. Stay consistent. Get results.”
2. Problem Statement
   Many young people, especially in Africa, struggle with:
   Staying consistent and accountable while working or applying for jobs.


Scattered job search efforts across different platforms.


Weak application materials (CVs, cover letters, proposals).


Loss of motivation due to lack of results or feedback.


Aiki solves this by combining accountability, AI career assistance, and job search management in one clean, engaging app.
3. Target Users
   Students and recent graduates actively job-hunting.


Freelancers using platforms like Upwork or Contra.


Professionals seeking accountability to stay productive.


4. MVP Objectives
   The MVP will focus on core engagement, job application organization, and AI-assisted career support, while keeping the experience lightweight and affordable.
5. Core Features (MVP Scope)
   A. Onboarding & Profile Setup
   Gamified onboarding (Duolingo-inspired).


Profile creation: name, profession, goals, CV upload (optional).


Welcome message introducing “Lock In” sessions.


B. Lock-In Productivity Sessions
Button to “Lock In” for the day (start focused work).


Timer and streak counter.


Motivational prompts and badges (gamification).


C. Job Application Tracker
Add jobs manually or from integrated job boards (Upwork, LinkedIn, Contra).


Columns: Job Title, Platform, Status (Applied / Interview / Offer).


Progress summary for each week.


D. AI Assistant (Career Support)
CV/Resume improvement suggestions.


Proposal and cover letter writer (short text prompts).


Quick feedback messages integrated in the app.


E. Notifications & Engagement
Daily reminder to “Lock In.”


Encouraging push notifications (“Keep your streak going!”).


Progress badges and weekly performance summary.


6. Design Requirements (UI/UX)
   Clean, minimal, and modern interface.


Light color palette (white, soft blue, mint green, or lavender accent).


Rounded corners, plenty of white space, and readable typography.


Smooth transitions, no clutter.


Duolingo-style gamification visuals (subtle, not cartoonish).


7. Success Metrics
   Daily active users (DAU)


Number of “Lock In” sessions per week


CV/AI feedback requests


Job applications tracked


8. Future Roadmap (Post-MVP)
   Full integration with job APIs (Upwork, LinkedIn).


Community leaderboard and group accountability rooms.


Advanced analytics for personal progress.


schema:
user [icon: user] {
id integer
first_name string
last_name string
email string
phone_number string
created_at string
}

user_profile [] {
id integer
user_id integer
user_link_id string
education_id integer
location string user address
}

user_skills [] {
id integer
user_id integer
name string
}

user_links [] {
id integer
user_id integer
upwork string
linkedin string
}

job_applications [icon: box] {
id integer
user_id integer
title string
platform string
status enum (APPLIED, INTERVIEWED, OFFERED, REJECTED)
created_at timestamp
updated_at timestamp
}

