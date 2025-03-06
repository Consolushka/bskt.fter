# IMP


## About
Service that evaluates player Impact on the final result (IMP). Based on:
<ul>
    <li>Minutes Played</li>
    <li>Plus-Minus</li>
    <li>Final Game Score</li>
</ul>

IMP provides deep analysis of player performance considering actual playing time volume.

## Supported Leagues

<ul>
    <li><a href="https://www.nba.com/">NBA</a></li>
    <li><a href="https://tambov.ilovebasket.ru/competitions/89960">Tambov MLBL</a></li>
</ul>

### CLI

<ul>
    <li>
        cron
        <p>
            Every day at 12AM service saves yesterday played games for each league.
        </p>
    </li>
    <li>
        seed:list
        <p>
            Prints all available database seeders.
        </p>
    </li>
    <li>
        seed:db
        <p>
            Seed database with provided seeder.
        </p>
    </li>
    <li>
        save-game
        <p>
            Save game from provided league with provided game id.
        </p>
    </li>
    <li>
        save-games-by-date
        <p>
            Save already played games for provided league and date (in format dd-mm-yyyy). If date is not provided, saves today games.
        </p>
    </li>
    <li>
        save-game-by-team
        <p>
            Save already played games for provided league and team id
        </p>
    </li>
    <li>
        serve
        <p>
            Starts API server. On <b>air</b>
        </p>
    </li>
</ul>

### API

Base path: <b>http://localhost:8001/api</b>

Swagger documentation: <b>http://localhost:8001/api/docs</b>