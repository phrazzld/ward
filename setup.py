from setuptools import setup

setup(
    name="ward-hooks",
    version="0.1.0",
    description="Pre-commit hooks for code review",
    author="Phaedrus Raznikov",
    author_email="phrazzld@pm.me",
    py_modules=[],
    scripts=[
        "ward_check.sh", 
        "ward_log.sh", 
        "ward_analyze.sh"
    ],
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
)