# Deployment Guide - AWS S3 Setup

## Step 1: Create S3 Bucket

### Option A: Using AWS CLI
```bash
# Create bucket (replace with your bucket name and region)
aws s3 mb s3://secure-file-sync-storage --region us-east-1

# Enable versioning (optional)
aws s3api put-bucket-versioning \
  --bucket secure-file-sync-storage \
  --versioning-configuration Status=Enabled

# Set bucket policy for security (optional)
aws s3api put-bucket-policy --bucket secure-file-sync-storage --policy file://bucket-policy.json
```

### Option B: Using AWS Console
1. Go to https://console.aws.amazon.com/s3/
2. Click "Create bucket"
3. Enter bucket name (e.g., `secure-file-sync-storage`)
4. Select region
5. Configure settings (versioning, encryption, etc.)
6. Create bucket

## Step 2: Create IAM User

### Using AWS CLI
```bash
# Create IAM user
aws iam create-user --user-name secure-file-sync-s3-user

# Create access key
aws iam create-access-key --user-name secure-file-sync-s3-user

# Attach policy for S3 access
aws iam put-user-policy \
  --user-name secure-file-sync-s3-user \
  --policy-name S3AccessPolicy \
  --policy-document file://s3-policy.json
```

### Using AWS Console
1. Go to https://console.aws.amazon.com/iam/
2. Click "Users" → "Create user"
3. Enter username: `secure-file-sync-s3-user`
4. Select "Attach policies directly"
5. Search and attach: `AmazonS3FullAccess` (or create custom policy)
6. Create user
7. Go to "Security credentials" tab
8. Click "Create access key"
9. Choose "Application running outside AWS"
10. Save the Access Key ID and Secret Access Key

## Step 3: IAM Policy (Custom - Recommended)

Create a file `s3-policy.json`:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:GetObject",
        "s3:DeleteObject"
      ],
      "Resource": "arn:aws:s3:::secure-file-sync-storage/*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:ListBucket"
      ],
      "Resource": "arn:aws:s3:::secure-file-sync-storage"
    }
  ]
}
```

## Step 4: Set Environment Variables on Fly.io

After obtaining your AWS Access Key ID and Secret Access Key:

```bash
# Set secrets on Fly.io
flyctl secrets set AWS_ACCESS_KEY_ID=your-access-key-id
flyctl secrets set AWS_SECRET_ACCESS_KEY=your-secret-access-key
flyctl secrets set AWS_REGION=us-east-1
flyctl secrets set AWS_S3_BUCKET=secure-file-sync-storage
```

## Step 5: Verify Secrets

```bash
# List secrets (values are hidden)
flyctl secrets list
```

## Notes

- Replace `secure-file-sync-storage` with your actual bucket name
- Replace `us-east-1` with your preferred AWS region
- Keep your AWS credentials secure - never commit them to git
- Use least-privilege IAM policies for production
- Consider using AWS Secrets Manager or Parameter Store for production

