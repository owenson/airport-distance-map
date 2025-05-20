import pandas as pd
import sqlite3

# Read the CSV file
print("Reading CSV file...")
df = pd.read_csv('airports.csv')

# Filter out heliports and closed airports
print("Filtering airports...")
df = df[~df['type'].isin(['heliport', 'closed'])]

# Filter for European airports (approximate boundaries)
df = df[
    (df['latitude_deg'] >= 35) & 
    (df['latitude_deg'] <= 72) & 
    (df['longitude_deg'] >= -10) & 
    (df['longitude_deg'] <= 40)
]

# Create SQLite connection
print("Creating SQLite database...")
conn = sqlite3.connect('airports.db')

# Write to SQLite
print("Writing to SQLite...")
df.to_sql('airports', conn, if_exists='replace', index=False)

# Create index on name column for faster searches
print("Creating index...")
conn.execute('CREATE INDEX IF NOT EXISTS idx_name ON airports(name)')

# Close connection
conn.close()
print("Done! Database created successfully.")
print(f"Total airports in database: {len(df)}") 