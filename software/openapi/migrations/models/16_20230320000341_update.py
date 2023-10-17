from tortoise import BaseDBAsyncClient


async def upgrade(db: BaseDBAsyncClient) -> str:
    return """
        ALTER TABLE "api" RENAME COLUMN "create_time" TO "created_time";
        ALTER TABLE "api" RENAME COLUMN "update_time" TO "updated_time";"""


async def downgrade(db: BaseDBAsyncClient) -> str:
    return """
        ALTER TABLE "api" RENAME COLUMN "created_time" TO "create_time";
        ALTER TABLE "api" RENAME COLUMN "updated_time" TO "update_time";"""
