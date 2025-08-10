@echo off
echo ========================================
echo Starting Admin Statistics API
echo ========================================

echo.
echo Starting MongoDB and API with Docker Compose...
docker-compose up -d

echo.
echo Waiting for services to be ready...
timeout /t 15 /nobreak > nul

echo.
echo API should now be running at http://localhost:8090
echo MongoDB should be running at localhost:27017
echo.

echo To test the API, run: test_api.bat
echo To stop the services, run: stop_api.bat
echo.

echo ========================================
echo Services Started!
echo ========================================
pause