from setuptools import setup, find_packages
import os

setup(
    name="event_helper",  # Package name
    version="0.1.0",
    packages=find_packages(),  # Automatically finds subpackages
    install_requires=[],  # Add dependencies here if needed
    author="UploadPilot",
    author_email="dev@uploadtpilot.com",
    description="A local package with event processing decorators.",
    long_description=open("README.md").read() if "README.md" in os.listdir() else "",
    long_description_content_type="text/markdown",
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
    ],
    python_requires=">=3.6",
)
